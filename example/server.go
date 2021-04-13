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
	egu.SecurePanic(app.LoadRouter(func(router *gin.Engine, logger *log.Logger) {
		// 引入跨域、recover、请求日志三个中间件
		router.Use(middleware.CORS, middleware.Recover, middleware.RequestLog(logger))
		router.GET("index", func(ctx *gin.Context) {
			// 记录日志
			logger.Info("test", log.Context{"hello": "world"})
			// 输出响应
			response.WrapContext(ctx).Success(nil)
		})
	}))
	// 启动http服务
	app.ServeHTTP()
	// 启动应用
	app.Run()
}
