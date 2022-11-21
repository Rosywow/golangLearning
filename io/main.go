package main

import (
	"bufio"
	"os"
)

func main()  {
	file, err := os.Open("./test")
	if err != nil{
		panic(err)
	}

	bufio.NewReader(file)
}
