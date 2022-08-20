package rebuild

import (
	"context"
	"github.com/ebar-go/ego/rebuild/application"
	"github.com/ebar-go/ego/rebuild/component"
	"github.com/ebar-go/ego/rebuild/server"
)

type Application interface {
	WithComponent(component ...component.Component) Application
	WithServer(server ...server.Server) Application
	Run(stop <-chan struct{}) error
}

func Run(options ServerRunOptions) {
	app := application.NewApplication()
	app.WithComponent(component.NewCache(), component.NewLogger())
	app.WithServer(server.NewHTTPServer(options.HttpAddr))

	component.Provider().Logger().Info("Application started")
	app.Run(context.Background().Done())
	component.Provider().Logger().Info("Application stopped")
}
