package main

import (
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/response"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 初始化应用
	app := ego.App()
	// 加载配置文件
	if err := app.LoadConfig("./app.yaml"); err != nil {
		log.Fatalf("failed to load config: %v\n", err)
	}

	// 初始化路由
	err := app.LoadRouter(func(router *gin.Engine) {
		// 引入跨域、recover、请求日志三个中间件
		router.Use(middleware.CORS, middleware.Recover)
		router.GET("index", func(ctx *gin.Context) {
			// 输出响应
			response.WrapContext(ctx).Success(nil)
		})
		router.Any("check", func(ctx *gin.Context) {
			// 输出响应
			response.WrapContext(ctx).Success(nil)
		})
	})
	if err != nil {
		log.Fatalf("failed to load router:%v\n", err)
	}
	// 启动http服务
	app.ServeHTTP()
	// 启动应用
	app.Run()
}
