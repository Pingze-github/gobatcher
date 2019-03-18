package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
func request(cli *http.Client, uri string) {
	// 复用client
	resp, err := cli.Get(uri)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("return" + strconv.Itoa(resp.StatusCode))
	}
}

func request2(cli *fasthttp.Client, uri string) {
	req := &fasthttp.Request{}
	req.SetRequestURI(uri)
	resp := &fasthttp.Response{}
	//	stime := time.Now()
	if err := cli.DoTimeout(req, resp, time.Second*30); err != nil {
		panic(errors.New("[Fasthttp Error]" + err.Error()))
	}
	if resp.StatusCode() != 200 {
		panic("return" + strconv.Itoa(resp.StatusCode()))
	}
	// dtime := time.Now().Sub(stime)
	// log.Println(resp.StatusCode(), dtime)
}

func workRoutine(ch chan string, wg *sync.WaitGroup) {
	mode := 2
	// cli.MaxConnsPerHost = 2048
	for uri := range ch {
		if mode == 1 {
			request(&http.Client{}, uri)
		} else {
			request2(&fasthttp.Client{}, uri)
		}
	}
	wg.Done()
}

// 如何提高服务能够产生的并发量
// 分时统计
// 参照jmeter
func main() {
	rnum := 128
	tnum := 128000

	divide := 1
	rnum = rnum / divide
	tnum = tnum / divide
	log.Println("并发数", rnum)

	url := "http://127.0.0.1:9001/"
	// url := "https://im.qq.com/pcqq/"

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
