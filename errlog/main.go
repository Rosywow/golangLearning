package main

import (
	"errlog/happy"
	"log"
	"os"
)

func main()  {
	log.Println("123")
	log.SetFlags(log.Ldate |log.Lshortfile |log.Ltime)
	file,err :=os.OpenFile("happy/logfile",os.O_APPEND|os.O_CREATE|os.O_RDWR,0777)
	if err!=nil{
		log.Fatal(err)
	}
	log.SetOutput(file)
	happy.Test()
}