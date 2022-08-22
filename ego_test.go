package ego

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	Run(ServerRunOptions{HttpAddr: ":8080", HttpTraceHeader: "trace", RPCAddr: ":8081", WebSocketAddr: ":8082"})
}

func Run(options ServerRunOptions) {
	aggregator := NewAggregatorServer()
	aggregator.WithComponent(component.NewCache(), component.NewLogger())

	// http server example
	httpServer := NewHTTPServer(options.HttpAddr).
		WithDefaultRecoverMiddleware().
		WithDefaultRequestLogMiddleware().
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
			router.GET("/panic", func(ctx *gin.Context) {
				if ctx.Query("id") == "" {
					panic(errors.InvalidParam("invalid params"))
				}
				ctx.String(http.StatusOK, "panic")
			})
		})

	// grpc server example
	grpcServer := NewGRPCServer(options.RPCAddr).
		WithKeepAliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     time.Minute,
			MaxConnectionAge:      time.Hour,
			MaxConnectionAgeGrace: time.Hour * 24,
			Time:                  time.Hour,
			Timeout:               time.Second * 30,
		}).
		RegisterService(func(s *grpc.Server) {
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
