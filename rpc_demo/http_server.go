package main

import (
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	http.HandleFunc("/panda", func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, "hello world")
	})
	l, err := net.Listen("tcp", ":10086")
	if err != nil {
		log.Fatalln("网络错误")
	}
	defer l.Close()
	http.Serve(l, nil)
}
