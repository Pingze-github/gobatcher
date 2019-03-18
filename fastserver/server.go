package main

import (
	"fmt"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

func handler(ctx *fasthttp.RequestCtx) {
	time.Sleep(time.Millisecond * 10)
	fmt.Fprintf(ctx, "")
}

func main() {
	log.Println("FastServer Start ...")
	fasthttp.ListenAndServe(":9001", handler)
}
