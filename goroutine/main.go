package main

import (
	"fmt"
	"time"
)

func main() {
	//var wg sync.WaitGroup


	//进行多次测试
	for  {
		//单次测试，持续时间为1小时
		t := time.NewTicker(2*time.Second)
		arr := []int{1,2,3,4}
		for i:=0;i<2;i++ {
			//每个协程一个小时后就会退出
			go func(ticker *time.Ticker) {
				FOR:
				for{
					select {
					case <-ticker.C:
						{
							//已超时，退出协程
							break FOR
						}
					default:
						{
							fmt.Println("123")
						}

					}
				}
			}(t)
		}

		//结束后close掉所有conn
		//结束所有go func
		time.Sleep(5*time.Second)
		fmt.Println("再停止3秒，进入下一个循环")
		fmt.Println(arr)
		arr = nil
		fmt.Println(arr)
		time.Sleep(3*time.Second)
	}
}
