package main

import (
	"fmt"
	"sync"
	"time"
)

var mutax = sync.Mutex{}

func main() {
	go func() {
		mutax.Lock()
		fmt.Println("协程1")
		time.Sleep(2*time.Second)
		mutax.Unlock()
	}()
	go func() {
		mutax.Lock()
		fmt.Println("协程2")
		time.Sleep(2*time.Second)
		mutax.Unlock()
	}()
	for  {
		
	}
}