package aggregator

import (
	"context"
	component "github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/runtime"
	"github.com/ebar-go/ego/server"
	"time"
)

// Aggregator includes the components and servers
type Aggregator struct {
	provider *component.Container
	servers  []server.Server
}

func (a *Aggregator) WithComponent(components ...component.Component) *Aggregator {
	a.provider.Register(components...)
	return a
}

func (a *Aggregator) WithServer(servers ...server.Server) *Aggregator {
	a.servers = append(a.servers, servers...)
	return a
}

func (a *Aggregator) prepare() {
	component.Initialize(a.provider)
}

func (a *Aggregator) Run() {

	ctx, cancel := context.WithCancel(context.Background())
	for _, s := range a.servers {
		go s.Serve(ctx.Done())
	}

	component.Provider().Logger().Info("Aggregator started successfully")

	runtime.Shutdown(func() {
		cancel()
		time.Sleep(time.Second)
		component.Provider().Logger().Info("Aggregator stopped successfully")
	})
}

func NewAggregator() *Aggregator {
	return &Aggregator{provider: component.NewContainer()}
}

type NamedAggregator struct {
	*Aggregator
	name string
}

func (a *NamedAggregator) Run() {
	component.Provider().Logger().Infof("Running Aggregator:%s\n", a.name)
	a.Aggregator.Run()
}

func NewNamedAggregator(name string) *NamedAggregator {
	return &NamedAggregator{name: name}
}
