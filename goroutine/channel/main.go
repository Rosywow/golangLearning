package main

import "fmt"

func main() {
	channel := make(chan interface{},1)
	channel <- 1
	close(channel)
	fmt.Println(<-channel)
}
