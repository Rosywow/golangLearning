package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main(){
	ctx := context.Background()
	ctx,cancel := context.WithTimeout(ctx,time.Second)
	defer cancel()
	req,err :=http.NewRequest(http.MethodGet,"http://localhost:8080",nil)
	if err!=nil {
		log.Fatal(err.Error())
	}
	req = req.WithContext(ctx)
	response,err :=http.DefaultClient.Do(req)
	if err!=nil {
		log.Fatal(err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode!=http.StatusOK {
		log.Fatal(response.Status)
	}

	io.Copy(os.Stdout,response.Body)

}
