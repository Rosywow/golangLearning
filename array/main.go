package main

import "fmt"

func main(){
	letters := []string{"a","b","c","d"}
	fmt.Println(len(letters))
	letters=nil
	fmt.Println(len(letters))
}
