package ego

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/server"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	Run(ServerRunOptions{HttpAddr: ":8080", HttpTraceHeader: "trace", RPCAddr: ":8081", WebSocketAddr: ":8082"})
}

func Run(options ServerRunOptions) {
	aggregator := NewAggregatorServer()
	aggregator.WithComponent(component.NewCache(), component.NewLogger())

	// http server example
	httpServer := NewHTTPServer(options.HttpAddr).
		EnableReleaseMode().
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

	// grpc server example
	grpcServer := NewGRPCServer(options.RPCAddr).RegisterService(func(s *grpc.Server) {
		// pb.RegisterGreeterServer(s, &HelloService{})
	})

	// websocket server example
	websocketServer := NewWebsocketServer(options.WebSocketAddr).OnConnect(func(conn server.Conn) {
		// connected
	}).OnDisconnect(func(conn server.Conn) {
		// disconnected
	}).OnMessage(func(conn server.Conn, msg []byte) {
		// on request
	})

	aggregator.WithServer(httpServer, grpcServer, websocketServer)

	aggregator.Run()
}
