package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// 例如，go build 生成可执行文件scanner.exe之后
	// 在命令行输入 scanner.exe 就可以运行该程序， 此时
	// os.Args = 1, 如果输入的是 scanner.exe main.go
	// 此时，os.Args = 2
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage:\n\t%s [files]\n", os.Args[0])
		os.Exit(1)
	}

	// Before we can call the Init method in scanner.
	// Scanner we will read the file contents and
	// create a token.FileSet holding a token.File
	// per file we scan.
	fs := token.NewFileSet()

	for _, arg := range os.Args[1:] {
		b,err := ioutil.ReadFile(arg)
		if err != nil {
			log.Fatal(err)
		}

		f := fs.AddFile(arg,fs.Base(),len(b))
		var s scanner.Scanner
		s.Init(f,b,nil,scanner.ScanComments)


		// Once the scanner has been initialized we
		// can call Scan and print the token we obtain.
		// Once we reach the end of the file scanned,
		// we will obtain an EOF (End Of File) token.
		for {
			_, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			fmt.Println(tok, lit)
		}
	}
}
