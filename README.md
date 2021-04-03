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

## 安装
```
go get -u github.com/ebar-go/ego
```

## 示例
- 配置文件
```yaml
server:
  name: demo
http:
  port: 8085
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
	egu.SecurePanic(app.Container().Invoke(routeLoader))
	// 启动http服务
	app.ServeHTTP()
	// 启动应用
	app.Run()
}
// routerLoader 加载路由
func routeLoader(router *gin.Engine,logger *log.Logger,)  {
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

- 通过`go run main.go`启动服务
```
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /index                    --> main.initRouter.func1 (6 handlers)
2021-03-29 00:23:01.786178 I | Listening and serving HTTP on :8085
```

访问`localhost:8085/index`验证结果。

## 文档
详细文档地址：[https://ebar-go.github.io](https://ebar-go.github.io)

