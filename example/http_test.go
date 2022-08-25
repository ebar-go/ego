package example

import (
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/component"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"testing"
	"time"
)

func TestAggregatorWithHttpServer(t *testing.T) {
	aggregator := ego.NewAggregatorServer()

	// 实例化一个http服务
	httpServer := ego.NewHTTPServer(":8080").
		RegisterRouteLoader(func(router *gin.Engine) { // 注册路由
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "home")
			})
		})

	aggregator.WithServer(httpServer)

	aggregator.Run()
}

func TestAggregatorWithHttpServerComplex(t *testing.T) {
	aggregator := ego.NewAggregatorServer()

	// 实例化一个http服务
	httpServer := ego.NewHTTPServer(":8080").
		WithDefaultRecoverMiddleware(). // 使用默认的recover组件
		WithDefaultRequestLogMiddleware(). // 使用默认的请求日志组件
		EnableReleaseMode(). // 开启release mode
		EnablePprofHandler(). // 开启pprof分析
		EnableAvailableHealthCheck(). // 开启健康检查
		EnableSwaggerHandler(). // 开启swagger接口插件
		EnableCorsMiddleware(). // 开启跨域组件
		EnableTraceMiddleware("traceHeader"). // 开启全局链路组件
		WithNotFoundHandler(func(ctx *gin.Context) { // 配置404
			ctx.String(http.StatusNotFound, "404 Not Found")
		}).
		RegisterRouteLoader(func(router *gin.Engine) { // 注册路由
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "home")
			})
		})

	aggregator.WithServer(httpServer)

	aggregator.Run()
}

func TestAggregatorWithComponents(t *testing.T) {
	aggregator := ego.NewAggregatorServer()

	// component logger
	component.Provider().Logger().Info("test logger info function")

	// component cache
	component.Provider().Cache().Default().Set("someCacheKey", "someCacheValue", time.Minute)

	// component redis
	if err := component.Provider().Redis().Connect(&redis.Options{
		// ... some options like address, port
	}); err != nil {
		panic(err)
	}
	component.Provider().Redis().Set("someRedisKey", "someRedisVal", time.Minute)

	// 实例化一个http服务
	httpServer := ego.NewHTTPServer(":8080").
		RegisterRouteLoader(func(router *gin.Engine) { // 注册路由
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "home")
			})
		})

	aggregator.WithServer(httpServer)

	aggregator.Run()
}
