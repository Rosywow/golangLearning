package Text

import (
	"github.com/gorilla/websocket"
	"log"
	"testing"
	"time"
)

func TestWebSocket(t *testing.T) {
	dl := websocket.Dialer{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
	//限制通道最大值为2000
	ch :=make(chan int,2000)
	sumTime := time.NewTicker(5*time.Minute)

FOR:
	for {
		//检查有没有超过测试时长
		select {
		case <-sumTime.C:
			{
				//已经超过总时长
				log.Println("测试结束")
				t.Logf("通过测试")
				break FOR
			}
		default:
			{
				//还没超过总时长，什么都不执行
			}

		}


		t1 := time.NewTicker(3*time.Second)
		ch <- 1
		go func (ticker *time.Ticker){
			select {
			case <-ticker.C:
				{
					//单个响应已经超时了
					log.Fatal("响应时间大于三秒")
					<-ch
					t.Error("单次响应时间大于三秒")
				}
			default:
				{
					//单个响应还没超时
					url := "ws://localhost:9090/api/chat?sender=1&receiver=2"
					c, _, err := dl.Dial(url, nil)
					if err != nil {
						log.Println("连接失败:", err)
					}else {
						//fmt.Printf("响应:%s", fmt.Sprint(res))
						c.Close()
					}
					<-ch
				}
			}
		}(t1)
	}
}
