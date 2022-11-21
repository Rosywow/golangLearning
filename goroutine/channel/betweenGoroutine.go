package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	//channel := make(chan int)
	//var wg sync.WaitGroup
	//wg.Add(2)
	//
	//go func(ch chan int) {
	//	<- ch
	//	fmt.Println("终于等到你")
	//	wg.Done()
	//}(channel)
	//
	//go func(ch chan int) {
	//	time.Sleep(3*time.Second)
	//	ch <- 1
	//	wg.Done()
	//}(channel)
	//
	//wg.Wait()
	//fmt.Println("finish")
	//user := 1
	//for {
	//	if user%2 == 1 {
	//		fmt.Println()
	//	}
	//	user += 1
	//}


	//两个子进程两两为一组的通信
	var wg sync.WaitGroup
	user := 1
	var channel chan int
	wg.Add(8)
	for i:=1;i<9;i++ {
		//单数
		if user%2 != 0 {
			fmt.Println(i)
			channel = make(chan int)
			go func(ch chan int) {
				fmt.Println("ch:",ch)
				fmt.Println("from: ",<-ch)
				wg.Done()
			}(channel)
		} else {
			fmt.Println(i)
			go func(ch chan int,num int) {
				time.Sleep(1*time.Second)
				ch <- num
				wg.Done()
			}(channel,i)
		}
		user++
	}
	wg.Wait()
	fmt.Println("finish")

}
