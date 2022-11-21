package Text

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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


func main()  {
	log.SetFlags(log.Ldate |log.Lshortfile )
	file,err :=os.OpenFile("Text/logFile",os.O_APPEND|os.O_CREATE|os.O_RDWR,0777)
	if err!=nil {
		log.Println("Open logfile err :",err)
	}
	log.SetOutput(file)
	dl := websocket.Dialer{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		//HandshakeTimeout: 3 * time.Second,
	}

	log.Println("开始测试")
	max := 512
	add := 0
	for {
		a,b := 1,2
		for b <= max+add {
			url :=fmt.Sprintf("ws://localhost:9090/api/chat?sender=%d&receiver=%d",a,b)
			c, res, err := dl.Dial(url, nil)
			if err!=nil {
				fmt.Println("connect err:",err,res)

			} else{
				Coon := ConnsType{
					Conns: c,
					Sender: fmt.Sprintf("%d",a),
					Receiver: fmt.Sprintf("%d",b),
				}
				testConns=append(testConns,Coon)
			}
			if a<b {
				temp:=a
				a=b
				b=temp
			}else {
				a+=1
				b+=3
			}
		}
		fmt.Println("连接建立完毕，倒数3秒")
		time.Sleep(3*time.Second)


		for _,v:=range testConns{
			go func(connect ConnsType){
			FOR:
				for{
					sendTime := time.Now()
					err = connect.Conns.WriteMessage(1,[]byte(fmt.Sprintf(`{"message":"你好","senderid":"%s","time":"%s"}`,connect.Sender,sendTime.Format("01-02-2006 15:04:05"))))
					if err!=nil{
						fmt.Println("write message err",err)
					}
					_,p,err :=connect.Conns.ReadMessage()
					if err!=nil {
						fmt.Println("1",err)
						break FOR
					}
					var Data MessageType
					err = json.Unmarshal(p,&Data)
					if err!=nil {
						fmt.Println("err",err)
					}

					//取出对方发送的消息的事件
					getSendTime,err := time.ParseInLocation("01-02-2006 15:04:05",Data.Time,time.Local)
					if err!=nil {
						fmt.Println("time.ParseInLocation err:",err)
					}
					different := time.Now().Sub(getSendTime)
					fmt.Println("different",different,Data.Message)
					//if different > 2*time.Second {
					//	fmt.Println("different:",different)
					//}
					//如果发送事件和接收事件超过x秒，则说明失败
					if different > 3*time.Second {
						log.Fatal("单次响应超时，此时并发数为:",max)
						return
					}
				}
			}(v)
		}
		time.Sleep(5*time.Minute)
		fmt.Println("一轮测试结束")
		for _,v:= range testConns{
			v.Conns.Close()
		}
		testConns = nil
		add+=200
		time.Sleep(10*time.Second)
	}
	////开2000个协程轮询
}

type LoginMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignIn(){
	done := make(chan bool,1)
	t := time.NewTicker(3*time.Second)
	go timer(t,done)

	url := "http://localhost:9090/api/signin"
	contentType := "application/json"
	requestBody := LoginMessage{ "1","123456789"}

	requestBodyJson,err := json.Marshal(requestBody)
	if err!=nil {
		fmt.Println(err)
		return
	}

	resp,err := http.Post(url,contentType,bytes.NewBuffer(requestBodyJson))
	if err!=nil {
		fmt.Println(err)
		return
	}
	data,_ := ioutil.ReadAll(resp.Body)
	fmt.Println(" Login resp",string(data))

	done <- true
}

func Login(){
	done := make(chan bool,1)
	t := time.NewTicker(3*time.Second)
	go timer(t,done)

	url := "http://localhost:9090/api/login"
	contentType := "application/json"
	requestBody := LoginMessage{ "1","123456789"}

	requestBodyJson,err := json.Marshal(requestBody)
	if err!=nil {
		fmt.Println(err)
		return
	}

	resp,err := http.Post(url,contentType,bytes.NewBuffer(requestBodyJson))
	if err!=nil {
		fmt.Println(err)
		return
	}

	data,_ := ioutil.ReadAll(resp.Body)
	fmt.Println("Login resp",string(data))

	done <- true
}

// UploadFile https://blog.csdn.net/huobo123/article/details/104288030
func UploadFile(){
	done := make(chan bool,1)

	t := time.NewTicker(3*time.Second)
	go timer(t,done)
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

	done <- true
}

func DownloadFile(){
	done := make(chan bool,1)
	t := time.NewTicker(3*time.Second)
	go timer(t,done)

	filepath := "main.go"
	url := fmt.Sprintf("http://localhost:1234/api/downloadfile?filepath=%s",filepath)

	resp,err := http.Get(url)
	if err!=nil {
		fmt.Println(err)
		return
	}
	SourseFile, err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		fmt.Println("ReadAll err: ",err)
		return
	}

	targetFile, err := os.Create("test.go")
	if err!=nil {
		fmt.Println("Open err: ",err)
		return
	}
	_, err = targetFile.Write(SourseFile)
	if err!=nil {
		fmt.Println("Write err: ",err)
	}

	done <- true
}

func timer(ticker *time.Ticker,Done chan bool) {

	select {
	case <-ticker.C:
		//超时了,在删除对应那一组conn，同时也删除那一组协程
		fmt.Println("超时了")
	case <-Done:
		//完成且，没超时
		fmt.Println("完成且，没超时")
	}
}

