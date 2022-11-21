package main

import (
	"fmt"
	"time"
)

func main() {
	//将time.Time转化为string
	now := time.Now().Format("01-02-2006 15:04:05")

	fmt.Println(now)
	//format := now.Format("01-02-2006 15:04:05 Monday")
	//fmt.Println(format)
	time.Sleep(3*time.Second)
	//end :=time.Now().Format("15:04:05")
	//fmt.Println("different:",end.Sub(now))

	//将字符串转化为time.Time类型
	first,err := time.ParseInLocation("01-02-2006 15:04:05",now,time.Local)
	if err!=nil {
		fmt.Println("time parse:",err)
	}
	fmt.Println(time.Now())


	//时间的减法
	different := time.Now().Sub(first)
	fmt.Println("different:",different)
}
