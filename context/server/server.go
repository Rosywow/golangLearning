package main

import (
	"context/mylog"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main(){
	http.HandleFunc("/",mylog.Decorate(handler))
	log.Fatal(http.ListenAndServe("127.0.0.1:8080",nil))
}
func handler(w http.ResponseWriter,r *http.Request){

	ctx := r.Context()

	mylog.Println(ctx,"handler started")
	defer mylog.Println(ctx,"handler stopped")

	select {
		case <-time.After(5*time.Second):
			fmt.Fprintf(w,"Hello World")
	    case <-ctx.Done():
	    	log.Print(ctx.Err())
	    	http.Error(w,ctx.Err().Error(),http.StatusInternalServerError)
	}

}