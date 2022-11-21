package main

import "fmt"

type data struct {
	A int
}

func main() {
	a := 100
	var Data data
	Data.A=a
	defer func() {
		a=12
		fmt.Println(Data.A)
	}()
}
