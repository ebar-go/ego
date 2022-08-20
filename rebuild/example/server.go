package main

import (
	"github.com/ebar-go/ego/rebuild/application"
	"github.com/ebar-go/ego/rebuild/component"
	"github.com/ebar-go/ego/rebuild/server"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
)

func main() {
	app := application.NewApplication()
	app.WithComponent(component.NewCache(), component.NewLogger())

	httpServer := server.NewHTTPServer(":8080").
		EnablePprofHandler().
		EnableAvailableHealthCheck().
		RegisterRouteLoader(func(router *gin.Engine) {
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "home")
			})
		})

	grpcServer := server.NewGRPCServer(":8081").RegisterService(func(s *grpc.Server) {
		// pb.pb.RegisterGreeterServer(s, &HelloService{})
	})

	app.WithServer(httpServer, grpcServer)

	app.Run()
}
