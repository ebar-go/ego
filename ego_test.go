package ego

import (
	"github.com/ebar-go/ego/component"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	Run(ServerRunOptions{HttpAddr: ":8080", HttpTraceHeader: "trace", RPCAddr: ":8081"})
}

func Run(options ServerRunOptions) {
	aggregator := NewAggregatorServer()
	aggregator.WithComponent(component.NewCache(), component.NewLogger())

	httpServer := NewHTTPServer(options.HttpAddr).
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

	grpcServer := NewGRPCServer(options.RPCAddr).RegisterService(func(s *grpc.Server) {
		// pb.RegisterGreeterServer(s, &HelloService{})
	})

	aggregator.WithServer(httpServer, grpcServer)

	aggregator.Run()
}
