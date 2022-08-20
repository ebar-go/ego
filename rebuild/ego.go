package rebuild

import (
	"github.com/ebar-go/ego/rebuild/application"
	"github.com/ebar-go/ego/rebuild/component"
	"github.com/ebar-go/ego/rebuild/server"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
)

type Application interface {
	WithComponent(component ...component.Component) Application
	WithServer(server ...server.Server) Application
	Run()
}

func Run(options ServerRunOptions) {
	app := application.NewApplication()
	app.WithComponent(component.NewCache(), component.NewLogger())

	httpServer := server.NewHTTPServer(options.HttpAddr).
		EnablePprofHandler().
		EnableAvailableHealthCheck().
		RegisterRouteLoader(func(router *gin.Engine) {
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "home")
			})
		})

	grpcServer := server.NewGRPCServer(options.RPCAddr).RegisterService(func(s *grpc.Server) {
		// pb.RegisterGreeterServer(s, &HelloService{})
	})

	app.WithServer(httpServer, grpcServer)

	app.Run()
}
