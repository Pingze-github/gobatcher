package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

func readFromResp(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// net/http/client 效率一般
func request() {
	// 复用client
	resp, err := http.Get("http://127.0.0.1:9000/")
	if err != nil {
		// Under heavy load with GOMAXPROCS>>1, it frequently fails
		// with transient failures like:
		// "dial tcp: cannot assign requested address"
		// or:
		// "ConnectEx tcp: Only one usage of each socket address
		// (protocol/network address/port) is normally permitted".
		// So we just log and continue,
		// otherwise significant fraction of benchmarks will fail.
		log.Printf("Get: %v", err)
		return
	}
	defer resp.Body.Close()
	// fmt.Println(resp.Status)
	// body, _ := readFromResp(resp)
	// fmt.Println(body)
}

func request2(cli *fasthttp.Client, uri string) {
	req := &fasthttp.Request{}
	req.SetRequestURI(uri)
	resp := &fasthttp.Response{}
	if err := cli.DoTimeout(req, resp, time.Second*10); err != nil {
		panic(errors.New("[Fasthttp Error]" + err.Error()))
	}
	// log.Println(resp.StatusCode())
}

func workRoutine(ch chan string, wg *sync.WaitGroup) {
	cli := &fasthttp.Client{}
	// cli.MaxConnsPerHost = 2048
	for uri := range ch {
		request2(cli, uri)
	}
	wg.Done()
}

// 如何提高服务能够产生的并发量
// 分时统计
// 参照jmeter
func main() {
	rnum := 2048
	tnum := 409600

	divide := 8
	rnum = rnum / divide
	tnum = tnum / divide
	log.Println("并发数", rnum)

	url := "http://127.0.0.1:9005/"

	ch := make(chan string)
	var wg sync.WaitGroup

	go func() {
		for i := 0; i < tnum; i++ {
			ch <- url
		}
		close(ch)
	}()

	wg.Add(rnum)
	for i := 0; i < rnum; i++ {
		go workRoutine(ch, &wg)
	}

	stime := time.Now()
	wg.Wait()
	dtime := time.Now().Sub(stime)
	fmt.Println("总耗时", dtime)
	fmt.Println("吞吐量", float64(tnum)/dtime.Seconds())
	fmt.Println("平均耗时", dtime.Seconds()/float64(tnum))
}
