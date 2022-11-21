package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func handleUpload(w http.ResponseWriter,r *http.Request){
	fmt.Println("laile")
	r.ParseMultipartForm(10 << 20)
	fmt.Println("r.MultipartForm",r.MultipartForm)
	fmt.Println("r.Form",r.Form)
	fmt.Println("r.Form",r.Body)
	fmt.Println("MD5",r.FormValue("MD5"))
	fmt.Println("filename",r.FormValue("filename"))
	//r.FormFile("file")
	w.Write([]byte(`{"status":200,"msg":"ok"}`))
}

func download(w http.ResponseWriter,r *http.Request){
	fmt.Println("downloadFile")
	filepath := r.URL.Query().Get("filepath")
	file,err := os.Open(filepath)
	if err!=nil {
		fmt.Println("Open err: ",err)
		return
	}

	filebytes,err := ioutil.ReadAll(file)
	if err!=nil{
		fmt.Println("ReadAll err: ",err)
		return
	}
	w.Write(filebytes)
}

func main() {
	http.HandleFunc("/api/upload",handleUpload)
	http.HandleFunc("/api/downloadfile",download)
	http.ListenAndServe(":1234",nil)
}
