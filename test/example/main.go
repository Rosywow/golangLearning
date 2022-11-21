package main

import (
"bytes"
"encoding/json"
"errors"
"fmt"
"github.com/gorilla/websocket"
"io"
"io/ioutil"
"log"
"mime/multipart"
"net/http"
"os"
"sync"
"sync/atomic"
"time"
)

var UP = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	HandshakeTimeout: 3 * time.Second,
}

type MessageType struct {
	Message string `json:"message"`
	Sender string `json:"sender"`
	Receiver string `json:"receiver"`
	Time string `json:"time"`
}
//
type ConnsType struct {
	Conns *websocket.Conn
	Sender string
	Receiver string
}

var testConns []ConnsType

type LoginMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	Status int `json:"status"`
}

func SignIn(username int) error {
	done := make(chan bool,1)
	msg := make(chan bool,1)
	t := time.NewTicker(3*time.Second)
	go timer(t,done,msg)

	url := "http://localhost:9090/api/signup"
	contentType := "application/json"
	usernameString := fmt.Sprintf("%d",username)
	requestBody := LoginMessage{ usernameString,"123456789"}

	requestBodyJson,err := json.Marshal(requestBody)
	if err!=nil {
		fmt.Println(err)
		return err
	}

	resp,err := http.Post(url,contentType,bytes.NewBuffer(requestBodyJson))
	if err!=nil {
		fmt.Println(err)
		return err
	}
	data,_ := ioutil.ReadAll(resp.Body)
	fmt.Println("signup resp",string(data))

	_ = Check(data)

	done <- true
	close(msg)
	x,_ := <-msg
	if x == true {
		return errors.New("overtime")
	}

	return nil
}

func Login(username int) error {
	done := make(chan bool,1)
	msg := make(chan bool,1)
	t := time.NewTicker(3*time.Second)
	go timer(t,done,msg)

	url := "http://localhost:9090/api/login"
	contentType := "application/json"
	usernameString := fmt.Sprintf("%d",username)
	requestBody := LoginMessage{ usernameString,"123456789"}

	requestBodyJson,err := json.Marshal(requestBody)
	if err!=nil {
		fmt.Println(err)
		return err
	}

	resp,err := http.Post(url,contentType,bytes.NewBuffer(requestBodyJson))
	if err!=nil {
		fmt.Println(err)
		return err
	}

	data,_ := ioutil.ReadAll(resp.Body)
	fmt.Println("Login resp",string(data))

	_ = Check(data)

	done <- true

	done <- true
	close(msg)
	x,_ := <-msg
	if x == true {
		return errors.New("overtime")
	}

	return nil
}

// UploadFile https://blog.csdn.net/huobo123/article/details/104288030
func UploadFile(){
	done := make(chan bool,1)
	msg := make(chan bool,1)
	t := time.NewTicker(3*time.Second)
	go timer(t,done,msg)
	file,err := os.Open("main.go")
	if err!=nil {
		fmt.Println("Open err: ",err)
		return
	}

	url := "http://localhost:1234/api/upload"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file","main.go")
	if err!=nil {
		fmt.Println("CreateFormFile err: ",err)
		return
	}

	_,err = io.Copy(part,file)
	if err!=nil {
		fmt.Println("io Copy err: ",err)
	}

	_ = writer.WriteField("MD5","qweqweqweqwe")
	_ = writer.WriteField("filename","main.go")

	err = writer.Close()
	if err!=nil {
		fmt.Println("writer Close err: ",err)
	}

	req, err := http.NewRequest("POST",url,body)
	if err!=nil {
		fmt.Println("NewRequest err: ",err)
	}
	req.Header.Add("Content-Type",writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err!=nil{
		fmt.Println("Do err: ",err)
	}
	data,_ := ioutil.ReadAll(resp.Body)
	fmt.Println("resp",string(data))

	err = Check(data)
	if err!=nil {
		fmt.Println("Check err:",err)
		return
	}

	done <- true
}

func DownloadFile(){
	done := make(chan bool,1)
	msg := make(chan bool,1)
	t := time.NewTicker(3*time.Second)
	go timer(t,done,msg)

	filepath := "main.go"
	url := fmt.Sprintf("http://localhost:1234/api/downloadfile?filepath=%s",filepath)

	resp,err := http.Get(url)
	if err!=nil {
		fmt.Println(err)
		return
	}
	SourceFile, err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		fmt.Println("ReadAll err: ",err)
		return
	}

	targetFile, err := os.Create("test.go")
	if err!=nil {
		fmt.Println("Open err: ",err)
		return
	}
	_, err = targetFile.Write(SourceFile)
	if err!=nil {
		fmt.Println("Write err: ",err)
	}

	done <- true
}

var mutex sync.Mutex
func Connect(sender int,receiver int)*websocket.Conn{
	dl := websocket.Dialer{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		//HandshakeTimeout: 3 * time.Second,
	}
	url :=fmt.Sprintf("ws://localhost:9090/api/chat?sender=%d&receiver=%d",sender,receiver)
	c, res, err := dl.Dial(url, nil)
	if err!=nil {
		fmt.Println("connect err:",err,res)
	} else{
		Coon := ConnsType{
			Conns: c,
			Sender: fmt.Sprintf("%d",sender),
			Receiver: fmt.Sprintf("%d",receiver),
		}
		mutex.Lock()
		testConns=append(testConns,Coon)
		mutex.Unlock()
	}
	return c
}

func sendMessage(conn *websocket.Conn,user int) {
	for  {
		sendTime := time.Now()
		err := conn.WriteMessage(1,[]byte(fmt.Sprintf(`{"message":"来自用户%d的消息","sender":"%d","time":"%s"}`,user,user,sendTime.Format("01-02-2006 15:04:05"))))
		if err!=nil{
			break
		}
		_, _, err = conn.ReadMessage()
		if err!=nil {
			break
		}
	}
}

func readMessage(conn *websocket.Conn,user int){
	for {
		var Data MessageType
		_, data, err := conn.ReadMessage()
		if err!=nil {
			break
		}

		err = json.Unmarshal(data,&Data)
		if err!=nil {
			fmt.Println("err",err)
		}
		getSendTime,err := time.ParseInLocation("01-02-2006 15:04:05",Data.Time,time.Local)
		if err!=nil {
			fmt.Println("time.ParseInLocation err:",err)
		}
		different := time.Now().Sub(getSendTime)
		fmt.Println("different: ",different)
		if different > 1*time.Second {
			log.Println("chat 超时")
			executeOnce.Do(changeStage)
			break
		}

		err = conn.WriteMessage(1,[]byte(fmt.Sprintf(`{"message":"用户%d接收到消息"}`,user)))
		if err!=nil{
			break
		}
	}
}

func Check(data []byte)error{
	var msg Message
	err := json.Unmarshal(data,&msg)
	if err!=nil{
		fmt.Println("Check err: ",err)
	}

	if msg.Status == 200 {
		//fmt.Println("成功")
		return nil
	}

	//fmt.Println("失败")
	return errors.New("err response")
}

func timer(ticker *time.Ticker,Done chan bool,Msg chan bool) {
	select {
	case <-ticker.C:
		//超时了,在删除对应那一组conn，同时也删除那一组协程
		executeOnce.Do(changeStage)
		fmt.Println("超时了")
		Msg <- true
	case <-Done:
		//完成且，没超时
		fmt.Println("完成,且没超时")
	}
}

// Once 让函数只执行一次
type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		// Outlined slow-path to allow inlining of the fast-path.
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

func (o *Once) recover() {
	o.m.Lock()
	defer o.m.Unlock()
	atomic.StoreUint32(&o.done, 0)
}

var executeOnce Once

func changeStage()  {
	increase = !increase
	fmt.Println("change stage to: ",increase)
}

var increase = true


func main(){
	var channel chan int
	user := 1
	for {
		if increase {
			if user%2 == 1 {
				channel = make(chan int)
			}
			go func(user int,isConnect chan int) {
				//err := SignIn(user)
				//if err!= nil {
				//	return
				//}
				//
				//err = Login(user)
				//if err!= nil {
				//	return
				//}
				//UploadFile()
				//DownloadFile()
				//单数发消息
				if user%2 == 1 {
					conn := Connect(user,user+1)
					//等待双数用户建立连接
					<-isConnect
					sendMessage(conn,user)
				} else {
					conn := Connect(user,user-1)
					//通知单数用户连接建立完毕，可以互发消息
					isConnect <- 1
					readMessage(conn,user)
				}
			}(user,channel)
			user += 1
		} else {
			for _,c := range testConns{
				_ = c.Conns.Close()
			}
			log.Println("当前并发数：",len(testConns))
			//清空数组
			testConns=nil
			executeOnce.recover()
			changeStage()
			fmt.Println(increase)
			//让用户总是从单数开始
			user += 1
			if user%2 == 0 {
				user += 1
			}
		}
		time.Sleep(1*time.Second)
	}
}


