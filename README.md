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
	"github.com/ebar-go/ego/component"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
)

func main() {
	Run(ego.ServerRunOptions{HttpAddr: ":8080", HttpTraceHeader: "trace", RPCAddr: ":8081"})
}

func Run(options ego.ServerRunOptions) {
	aggregator := ego.NewAggregatorServer()
	aggregator.WithComponent(component.NewCache(), component.NewLogger())

	httpServer := ego.NewHTTPServer(options.HttpAddr).
		EnablePprofHandler().
		EnableAvailableHealthCheck().
		EnableSwaggerHandler().
		EnableCorsMiddleware().
		EnableTraceMiddleware(options.HttpTraceHeader).
		WithNotFoundHandler(func(ctx *gin.Context) {
			ctx.String(http.StatusNotFound, "404 Not Found")
		}).
		RegisterRouteLoader(func(router *gin.Engine) {
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "home")
			})
		})

	grpcServer := ego.NewGRPCServer(options.RPCAddr).RegisterService(func(s *grpc.Server) {
		// pb.RegisterGreeterServer(s, &HelloService{})
	})

	aggregator.WithServer(httpServer, grpcServer)

	aggregator.Run()
}

```

## 文档
详细文档地址：[https://ebar-go.github.io](https://ebar-go.github.io)

