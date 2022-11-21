package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Msg struct {
	SenderId    int    `json:"sender"`
	ReceiverIds []int  `json:"receiverIds"`
	Message     string `json:"message"`
}

// 模拟2000个并发给服务器的不同接口发送请求，当某一个接口响应时间大于3秒时，结束测试
func main() {
//start:
	uid := 10000
	num := 300

	ctx, cancel := context.WithCancel(context.Background())
	doneCh := make(chan struct{})

	for i := uid; i <= uid+2000; i++ {
		go doTest(ctx, i, doneCh)
	}
	uid += 2000

	for {
		select {
		case <-time.After(10 * time.Second):
			for i := uid; i <= uid+num; i++ {
				go doTest(ctx, uid, doneCh)
			}
			uid += num

		case <-doneCh:
			cancel()
			time.Sleep(15 * time.Second)

			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			fmt.Printf("性能测试结束，超时并发数:%d ~ %d  内存使用情况:%d\n", uid-10000-num, uid-10000, m.Sys)
			//deleteInfo(num)
			//goto start
			return
		}
	}

}

func doTest(ctx context.Context, userId int, doneCh chan struct{}) {
	//register(userId, doneCh)
	//time.Sleep(2 * time.Second)

	go doWebsocket(userId, doneCh)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			login(userId, doneCh)
			time.Sleep(3 * time.Second)

			upload(doneCh)
			time.Sleep(3 * time.Second)

			download(doneCh)
			time.Sleep(3 * time.Second)
		}

	}

}

func download(doneCh chan struct{}) {
	connCh := make(chan string)

	url := "http://192.168.137.187:9090/api/file/download"
	//path := "D:\\Project\\CST\\myapp\\backend\\test\\testFile.txt.txt"
	param := fmt.Sprintf(`{"URL":"D:\\Project\\CST\\myapp\\backend\\test\\testFile.txt"}`)
	//fmt.Println("download params:" + param)

	go handleWithTimer(connCh, doneCh)
	body, err := httpPost(url, param)
	connCh <- "download"
	if err != nil {
		fmt.Println("download file error" + err.Error())
		return
	}

	fmt.Println(string(body))
	return
}

func upload(doneCh chan struct{}) {
	connCh := make(chan string)

	url := "http://192.168.137.187:9090/api/file/upload"
	path := "D:\\专题研究\\golang\\test\\testFile.txt"

	paramName := "uploadFile"

	req, err := newFileUploadRequest(url, paramName, path)
	if err != nil {
		return
	}
	go handleWithTimer(connCh, doneCh)
	client := &http.Client{}
	resp, err := client.Do(req)
	connCh <- "upload"
	if resp == nil {
		return
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(body))
	return

}

func login(userId int, doneCh chan struct{}) {
	connCh := make(chan string)

	url := "http://192.168.137.187:9090/api/user/login"
	phone := RandomPhone()
	params := fmt.Sprintf(`{"id":%d,"phone":"%s","password":"123456"}`, userId, phone)

	go handleWithTimer(connCh, doneCh)
	body, err := httpPost(url, params)
	if err != nil {
		fmt.Println("register error:" + err.Error())
		return
	}
	connCh <- "login"

	fmt.Println(string(body))

	return
}

func register(userId int, doneCh chan struct{}) {

	connCh := make(chan string)
	go handleWithTimer(connCh, doneCh)

	url := "http://192.168.137.187:9090/api/user/register"
	phone := RandomPhone()
	params := fmt.Sprintf(`{"id":%d,"phone":"%s","password":"123456"}`, userId, phone)

	_, err := httpPost(url, params)
	connCh <- "register"
	if err != nil {
		fmt.Println("register error:" + err.Error())
		return
	}
	return
}

func doWebsocket(userId int, done chan struct{}) {
	//conn := make(chan struct{})
	//handleWithTimer(conn, done)
	wsUrl := "ws://192.168.137.187:9090/api/ws?sender=%d"
	dialer := websocket.Dialer{}
	//向服务器发送连接请求，websocket 统一使用 ws://，默认端口和http一样都是80
	connect, _, err := dialer.Dial(fmt.Sprintf(wsUrl, userId), nil)
	if nil != err {
		log.Println(err)
	}
	//离开作用域关闭连接，go 的常规操作
	defer connect.Close()

	//定时向客户端发送数据
	go tickWriter(connect, userId)

	//启动数据读取循环，读取客户端发送来的数据
	for {
		//从 websocket 中读取数据
		//messageType 消息类型，websocket 标准
		//messageData 消息数据
		messageType, messageData, err := connect.ReadMessage()
		if nil != err {
			log.Println(err)
			break
		}
		switch messageType {
		case websocket.TextMessage: //文本数据
			fmt.Println(string(messageData))
			s := strings.Split(string(messageData), ":")

			t := s[4]

			t = t[:len(t)-2]
			k, _ := strconv.Atoi(t)
			spend := time.Now().UnixMilli() - int64(k)
			fmt.Printf("websocket spend %d ms\n", spend)

			if spend/1000 > 3 {
				fmt.Println("websocket timeout")
				done <- struct{}{}
			}
		case websocket.BinaryMessage: //二进制数据
			fmt.Println(messageData)
		case websocket.CloseMessage: //关闭
		case websocket.PingMessage: //Ping
		case websocket.PongMessage: //Pong
		default:
		}

	}
}

func tickWriter(connect *websocket.Conn, userId int) {
	for {
		// 1 and 2, 3 and 4
		receiver := userId + 1
		if userId%2 == 0 {
			receiver = userId - 1
		}
		s := fmt.Sprintf("hello user %d,i am %d,time:%d", receiver, userId, time.Now().UnixMilli())
		msg := Msg{SenderId: userId, ReceiverIds: []int{receiver}, Message: s}

		bytes, err := json.Marshal(msg)
		if err != nil {
			log.Println(err.Error())
		}

		//向客户端发送类型为文本的数据
		err = connect.WriteMessage(websocket.TextMessage, bytes)
		if nil != err {
			log.Println(err)
			break
		}

		time.Sleep(time.Second * 3)
	}
}

func deleteInfo(num int) {
	//请求结束后删除增加的用户
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 10000; i <= num+10000; i++ {

		i := i
		go func() {
			defer func() {
				wg.Done()
			}()

			url := "http://192.168.137.187:9090/api/admin/user-mgr"
			params := fmt.Sprintf(`{"id":%d}`, i)
			_, err := httpDelete(url, params)
			if err != nil {
				log.Printf("delete user %d error:%v\n", i, err.Error())
				return
			}
		}()
	}
	wg.Wait()
}

func handleWithTimer(connCh chan string, doneCh chan struct{}) {
	timer := time.NewTimer(6 * time.Second)
	start := time.Now().UnixMilli()

	select {
	case s := <-connCh:
		//获取到返回结果
		end := time.Now().UnixMilli()
		spend := end - start
		fmt.Printf("%s spend  %d ms\n", s, spend)
		return
	case <-timer.C: //超时
		fmt.Println("WaitChannel timeout")
		doneCh <- struct{}{}
		return
	}
}

func httpPost(url, params string) ([]byte, error) {
	var jsonStr = []byte(params)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if resp == nil || err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return nil, err
	}

	return body, nil
}

func httpDelete(url, params string) ([]byte, error) {
	var jsonStr = []byte(params)

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp == nil {
		return nil, err
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return nil, err
	}

	fmt.Println(string(body))
	return body, nil
}

func newFileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, path)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	//for key, val := range params {
	//	_ = writer.WriteField(key, val)
	//}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", uri, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request, err
}

const letterBytes = "0123456789"

var src = rand.NewSource(time.Now().UnixNano())

var headerNums = [...]string{"139", "138", "137", "136", "135", "134", "159", "158", "157", "150", "151", "152", "188", "187", "182", "183", "184", "178", "130", "131", "132", "156", "155", "186", "185", "176", "133", "153", "189", "180", "181", "177"}
var headerNumsLen = len(headerNums)

const (
	//使用四位二进制即可随机选择letterBytes里面的一位
	letterIdxBits = 4
	//掩码，即4个1
	letterIdxMask = 1<<letterIdxBits - 1
)
const (
	//使用六位二进制即可随机选择headerNums里面的一位
	headerIdxBits = 6
	headerIdxMask = 1<<headerIdxBits - 1
)

func getHeaderIdx(cache int64) int {
	for cache > 0 {
		//得到掩码对应的位数，比如1010110101 & 111111 = 0000110101
		//这样就可以取出二进制去随机选择数字了
		idx := int(cache & headerIdxMask)
		if idx < headerNumsLen {
			return idx
		}
		//取到的数字超过headerNums的长度，除去后6位，重新选择数字
		cache >>= headerIdxBits
	}
	//否则使用库函数生成
	return rand.Intn(headerNumsLen)
}

func RandomPhone() string {
	//12位手机号码
	b := make([]byte, 12)
	//获取一个63位的随机数
	cache := src.Int63()
	//获取选择手机号码前3位的随机数
	headerIdx := getHeaderIdx(cache)
	//使用得到的随机数去获取手机号码前3位
	for i := 0; i < 3; i++ {
		b[i] = headerNums[headerIdx][i]
	}
	//继续选择剩下的12位
	for i := 3; i < 12; {
		//生成的随机数用完了，重新生成
		if cache == 0 {
			cache = src.Int63()
		}
		//和getHeaderIdx一样
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
		cache >>= letterIdxBits
	}
	return string(b)
}