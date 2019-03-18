# GoBatcher
高并发压力测试。

# 问题

## 异常1
```
panic: Get http://127.0.0.1:9000/: dial tcp 127.0.0.1:9000: connectex: Only one usage of each socket address (protocol/network address/port) is normally permitted.
```

1. 复用http.Client
每个协程启动一个client，复用之。
问题依然存在。

2. 改用fasthttp
可行，整体效率提升近一倍。
3线程，吞吐量10K+
10线程，吞吐量11K+
100线程，吞吐量11K+
线程提升，对整体提升不明显。

但当服务器时延增大时，提升就很明显了。

513并发时，异常: `no free connections available to host`。这是fasthttp的`const DefaultMaxConnsPerHost = 512`设置引起的。
（似乎使用cli就没有限制了）
可以设置以突破限制（但吞吐量不会正比上升）：
```
cli := &fasthttp.Client{}
cli.MaxConnsPerHost = 1024
err := cli.DoTimeout(req, resp, timeout)
```

测试，net/http的server，空返回，最大平均吞吐量37K+。
在高并发数下，返回时间不影响吞吐量。
占用0%的CPU，内存占用1M。

fasthttp提供的server，空返回，最大并发平均吞吐量45K+（更大的线程数）。
占用0%的CPU，内存占用1M。

node的http.server，同样条件，平均吞吐量35K+。
占用25%的CPU，内存占用58M。

python，tornaodo异步server，最大吞吐量14K+。
不能支撑512线程。（最多400+）

[Windows]
Nginx反向代理fasthttp，默认配置下，最大吞吐量6K+。
不能支撑512线程。
+ worker_processes 1;
+ worker_connections 1024;
尝试更改配置也无法提高性能。


综合来看：
+ python表现很差，即使是tornado。
+ nodejs和go标准库的表现近似，都比较强。
+ golang使用fasthttp最强，达到45K+的吞吐量，同时资源消耗极小。