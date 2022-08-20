package rebuild

import (
	"context"
	"github.com/ebar-go/ego/rebuild/application"
	"github.com/ebar-go/ego/rebuild/component"
	"github.com/ebar-go/ego/rebuild/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Application interface {
	WithComponent(component ...component.Component) Application
	WithServer(server ...server.Server) Application
	Run(stop <-chan struct{}) error
}

func Run(options ServerRunOptions) {
	app := application.NewApplication()
	app.WithComponent(component.NewCache(), component.NewLogger())

	httpServer := server.NewHTTPServer(options.HttpAddr).
		EnableAvailableHealthCheck().
		RegisterRouteLoader(func(router *gin.Engine) {
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "home")
			})
		})
	app.WithServer(httpServer)

	component.Provider().Logger().Info("Application started")
	app.Run(context.Background().Done())
	component.Provider().Logger().Info("Application stopped")
}
