package application

import (
	"github.com/ebar-go/ego/rebuild/component"
	"github.com/ebar-go/ego/rebuild/server"
)

type Application struct {
	provider *component.Container
	servers  []server.Server
}

func (a *Application) WithComponent(components ...component.Component) *Application {
	a.provider.Register(components...)
	return a
}

func (a *Application) WithServer(server ...server.Server) *Application {
	a.servers = append(a.servers, server...)
	return a
}

type NamedApplication struct {
	*Application
	name string
}

func NewApplication() *Application {
	return &Application{provider: component.NewContainer()}
}

func (a *Application) prepare() {
	component.Initialize(a.provider)
}

func (a *Application) Run(stop <-chan struct{}) error {
	for _, s := range a.servers {
		go s.Serve(stop)
	}
	<-stop
	return nil
}
func NewNamedApplication(name string) *NamedApplication {
	return &NamedApplication{name: name}
}

func (a *NamedApplication) Run(stop <-chan struct{}) error {
	component.Provider().Logger().Infof("Running Application:%s\n", a.name)
	return a.Application.Run(stop)
}
