package ego

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/runtime"
	"github.com/ebar-go/ego/server"
	"github.com/ebar-go/ego/utils/async"
	"time"
)

// Engine includes the components and servers
type Engine struct {
	provider *component.Container
	servers  []server.Server
}

// WithComponent use the provided components.
func (engine *Engine) WithComponent(components ...component.Component) *Engine {
	engine.provider.Register(components...)
	return engine
}

// WithServer use the provided servers.
func (engine *Engine) WithServer(servers ...server.Server) *Engine {
	engine.servers = append(engine.servers, servers...)
	return engine
}

func (engine *Engine) prepare() {
	component.Initialize(engine.provider)
}

// BlockRun runs the engine with block until os.Exit.
func (engine *Engine) BlockRun() {
	engine.Run()
}

// Run runs the engine with non-blocking mode.
func (engine *Engine) Run() {
	go engine.Run()
}

func (engine *Engine) run() {
	engine.prepare()

	runner := async.NewRunner()
	for _, s := range engine.servers {
		runner.Add(func(stop chan struct{}) {
			s.Serve(stop)
		})
	}

	runner.Start()

	component.Provider().Logger().Info("Engine started successfully")

	runtime.Shutdown(func() {
		runner.Stop()
		time.Sleep(time.Second)
		component.Provider().Logger().Info("Engine stopped successfully")
	})
}

func buildEngine() *Engine {
	return &Engine{provider: component.NewContainer()}
}

// NamedEngine define engine with name.
type NamedEngine struct {
	*Engine
	name string
}

// Run override engine Run function.
func (engine *NamedEngine) Run() {
	component.Provider().Logger().Infof("Running Engine:%s\n", engine.name)
	engine.Engine.Run()
}

// NewNamedEngine creates a new named engine.
func NewNamedEngine(name string) *NamedEngine {
	return &NamedEngine{
		name:   name,
		Engine: buildEngine(),
	}
}
