package main

import (
	"fmt"
	"sync"
	"time"
)

func main(){
	var wg sync.WaitGroup
	i := 0
	for  {
		wg.Add(1)
		go func(sequence int) {
			fmt.Printf("我是用户%d 等待3秒增加一个新用户\n",sequence)
			time.Sleep(3*time.Second)
			wg.Done()
		}(i)
		wg.Wait()
		i++
	}
}
