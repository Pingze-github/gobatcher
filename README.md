# GoBatcher
高并发压力测试工具。

# 注意点
+ 应该在多个goroutine中，复用同个client。
+ 使用`github.com/valyala/fasthttp`代替`net/http`。

# 测试案例

> 本机请求，统一等待10ms，然后返回空字符串。

> 测试机：Windows10 Home | i5-7400 | 8G | Golang v1.11

## net/http server
+ 最大平均吞吐量 37K+
+ CPU 0%
+ 内存 ~1M
+ 最大并发 >1024

## fasthttp server
+ 最大平均吞吐量 45K+
+ CPU 0%
+ 内存 ~1M
+ 最大并发 >1024

## nodejs server
+ 最大平均吞吐量 35K+
+ CPU 25%
+ 内存 ~58M
+ 最大并发 >1024

## python tornado 异步server
+ 最大平均吞吐量 14K+
+ 最大并发 <512

## Nginx 反向代理 fasthttp server
+ 最大平均吞吐量 6K+
+ 最大并发 <512

> 测试机2：Debain8.9 | E5-2643 v4 | 8G | Golang v1.11

## fasthttp server
+ 最大平均吞吐量 37K+
+ CPU 0%
+ 内存 ~1M
+ 最大并发 >4096

## Nginx 反向代理 fasthttp server 优化配置
+ 最大平均吞吐量 14K+
+ 最大并发 <2048

## Nginx 直接return空字符串
+ 最大平均吞吐量 53K+
+ 最大并发 <4096

综合来看：
+ python表现很差，即使是tornado。
+ nodejs和go标准库的表现近似，都比较强。
+ golang使用fasthttp最强，达到45K+的吞吐量，同时资源消耗极小。
+ Nginx服务器的吞吐量很高，但并发数过高时，会产生错误`dialing to the given TCP address timed out`。
+ Nginx反向代理Golang服务时，总体性能会衰减许多。