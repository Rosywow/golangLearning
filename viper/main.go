package main

import (
	"fmt"
	"viper/config"
)

func main(){
	fmt.Println("apple:",config.Apple)
	config,err := config.GetViperValue()
	if err!=nil {
		fmt.Println("err:",err)
	} else {
		fmt.Println("config:",config)
	}
	username := config.DB.ContainerDB.User
	password := config.DB.ContainerDB.Password
	ip := config.DB.ContainerDB.IpAddress
	port := config.DB.ContainerDB.Port
	dbname := config.DB.ContainerDB.Dbname
	fmt.Println("username,password,ip,port,dbname",username,password,ip,port,dbname)
	s :=fmt.Sprintf("postgres://%s:%s@%s:%s/%s",username,password,ip,port,dbname)

	if s == "postgres://postgres:123456789@postgres:5432/postgres" {
		fmt.Println("相同")
	} else {
		fmt.Println("不相同")
	}
}
