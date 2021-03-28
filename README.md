## 说明
简单易用又强大的微服务golang框架。

## 特性
- http,websocket服务
- 丰富的中间件：请求日志、JWT认证,跨域,Recover,全局链路
- 全局服务:日志管理器，Redis,Mysql等
- 配置项
- 参数验证器
- curl组件
- Swagger

## 安装
```
go get -u github.com/ebar-go/ego
```

## 示例
- 配置文件
```yaml
server:
  name: demo
  httpPort: 8085
```

- main.go
```go
package main

import (
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/egu"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化应用
	app := ego.App()

	// 加载配置文件
	egu.SecurePanic(app.LoadConfig("./app.yaml"))

	// 初始化路由
	egu.SecurePanic(app.Container().Invoke(initRouter))

	// 启动http服务
	app.ServeHTTP()

	// 启动应用
	app.Run()
}

func initRouter(router *gin.Engine,logger *log.Logger,)  {
	// 引入跨域、recover、请求日志三个中间件
	router.Use(middleware.CORS, middleware.Recover, middleware.RequestLog(logger))
	router.GET("index", func(ctx *gin.Context) {
		// 记录日志
		logger.Info("test", log.Context{"hello":"world"})
		// 输出响应
		response.WrapContext(ctx).Success(nil)
	})
}
```
## 文档
详细文档地址：[https://ebar-go.github.io](https://ebar-go.github.io)

