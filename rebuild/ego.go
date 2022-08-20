package rebuild

import (
	"context"
	"github.com/ebar-go/ego/rebuild/component"
	"github.com/ebar-go/ego/rebuild/server"
)

type Schema struct {
	Protocol string
	Host     string
	Port     int
}

// Component represents a component
type Component interface{}

type Server interface {
	Serve(stop <-chan struct{}) error
}
type Application interface {
	WithComponent(component ...Component) Application
	WithServer(server ...Server) Application
	Run(stop <-chan struct{}) error
}

func NewApplication() Application {
	return nil
}

func Run() {
	application := NewApplication()
	application.WithComponent(component.NewCache(), component.NewLogger())
	application.WithServer(server.NewHTTPServer(""))

	component.Provider().Logger().Info("Application started")
	application.Run(context.Background().Done())
	component.Provider().Logger().Info("Application stopped")
}
