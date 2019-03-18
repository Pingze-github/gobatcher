package main

import (
	"log"
	"net/http"
	"time"
)

type handler struct{}

// 空返回
func (h *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Millisecond * 10)
	rw.Write([]byte(""))
}

func handlerFunc(rw http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Millisecond * 10)
	rw.Write([]byte(""))
}

func main() {
	log.Println("Server Start ...")

	// 最大吞吐20K
	// h := &handler{}
	// if err := http.ListenAndServe(":9000", &handler{}); err != nil {
	// 	panic(err)
	// }

	// 这种写法，最大吞吐38K
	http.HandleFunc("/", handlerFunc)
	if err := http.ListenAndServe(":9000", nil); err != nil {
		panic(err)
	}
}
