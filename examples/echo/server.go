package main

import (
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	aggregator := ego.New()

	aggregator.WithServer(httpServer())

	aggregator.Run()
}

func httpServer() server.Server {
	return ego.NewHTTPServer(":8080").EnableAvailableHealthCheck().
		WithDefaultRequestLogMiddleware().RegisterRouteLoader(func(router *gin.Engine) {
		router.GET("echo", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "Hello world")
		})
	})
}
