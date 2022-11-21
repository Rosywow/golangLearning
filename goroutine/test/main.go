package main

import "fmt"

func main() {
	arr := []string{"tom","james"}
	fmt.Println(len(arr),arr)
	arr = nil
	fmt.Println(len(arr),arr)
}
