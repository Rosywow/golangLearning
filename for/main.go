package main

import (
	"fmt"
	"log"
)

//看return在循环里面是否可以退出主程序
func main(){
	//for i:=0;i<5;i++ {
	//	fmt.Println("i:",i)
	//	for j:=0;j<5;j++ {
	//		fmt.Println("j:",j)
	//		if j==4 {
	//			return
	//		}
	//	}
	//}
	ch :=make(chan int,1)
	for i:=0;i<10;i++ {
		go func(num int) {
			if num==5{
				log.Println("退出主程序")
				ch<-1
			}
			for  {
				fmt.Println("i:",num)
			}
		}(i)
	}


	<-ch
}
