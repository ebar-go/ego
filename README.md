[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]

## 说明
简单易用又强大的微服务golang框架。

## 特性
- http服务
- 定时任务
- 丰富的中间件：请求日志、JWT认证,跨域,Recover,全局链路
- 集成Redis,Mysql,Jwt,Etcd客户端等基础组件
- 配置项
- 参数验证器
- curl组件
- Swagger

## Getting Started
- Install
```
go get github.com/ebar-go/ego
```

- main
```go
package main

import (
	"github.com/ebar-go/ego"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	aggregator := ego.NewAggregatorServer()

	httpServer := ego.NewHTTPServer(":8080").
		RegisterRouteLoader(func(router *gin.Engine) {
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "home")
			})
		})


	aggregator.WithServer(httpServer)

	aggregator.Run()
}

```

## 文档
详细文档地址：[https://ebar-go.github.io](https://ebar-go.github.io)

