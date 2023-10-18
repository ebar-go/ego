package main

import (
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/ebar-go/ego/utils/runtime/signal"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	aggregator := ego.New()

	aggregator.Aggregate(httpServer())

	aggregator.Run(signal.SetupSignalHandler())
}

func httpServer() runtime.Runnable {
	return ego.NewHTTPServer(":8090").
		EnableAvailableHealthCheck().
		EnablePprofHandler().
		WithDefaultRequestLogMiddleware().
		RegisterRouteLoader(func(router *gin.Engine) {
			router.GET("echo", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "Hello world")
			})
		})
}
