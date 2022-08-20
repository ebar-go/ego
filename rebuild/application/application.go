package application

import (
	"context"
	"github.com/ebar-go/ego/rebuild/component"
	"github.com/ebar-go/ego/rebuild/runtime"
	"github.com/ebar-go/ego/rebuild/server"
	"time"
)

type Application struct {
	provider *component.Container
	servers  []server.Server
}

func (a *Application) WithComponent(components ...component.Component) *Application {
	a.provider.Register(components...)
	return a
}

func (a *Application) WithServer(servers ...server.Server) *Application {
	a.servers = append(a.servers, servers...)
	return a
}

func NewApplication() *Application {
	return &Application{provider: component.NewContainer()}
}

func (a *Application) prepare() {
	component.Initialize(a.provider)
}

func (a *Application) Run() {

	ctx, cancel := context.WithCancel(context.Background())
	for _, s := range a.servers {
		go s.Serve(ctx.Done())
	}

	component.Provider().Logger().Info("Application started successfully")

	runtime.Shutdown(func() {
		cancel()
		time.Sleep(time.Second)
		component.Provider().Logger().Info("Application stopped successfully")
	})
}

type NamedApplication struct {
	*Application
	name string
}

func (a *NamedApplication) Run() {
	component.Provider().Logger().Infof("Running Application:%s\n", a.name)
	a.Application.Run()
}

func NewNamedApplication(name string) *NamedApplication {
	return &NamedApplication{name: name}
}
