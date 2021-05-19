package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	<-req.Context().Done()
	io.WriteString(w, "hello, world!\n")
}

func HelloServer2(w http.ResponseWriter, req *http.Request) {
	time.Sleep(10 * time.Second)
	io.WriteString(w, "hello, world!\n")
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	http.HandleFunc("/hello2", HelloServer2)
	err := http.ListenAndServe(":9002", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
