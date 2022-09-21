package ego

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/server"
	"github.com/ebar-go/ego/utils/async"
	"github.com/ebar-go/ego/utils/runtime"
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

// NonBlockingRun runs the engine with block until os.Exit.
func (engine *Engine) NonBlockingRun() {
	runtime.Goroutine(engine.run)
}

// Run runs the engine with blocking mode.
func (engine *Engine) Run() {
	engine.run()
}

func (engine *Engine) run() {
	engine.prepare()

	runner := async.NewRunner()
	for _, s := range engine.servers {
		runner.Add(s.Serve)
	}

	runner.Start()

	runtime.Shutdown(func() {
		runner.Stop()
		time.Sleep(time.Second)
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
	engine.Engine.Run()
}

// NewNamedEngine creates a new named engine.
func NewNamedEngine(name string) *NamedEngine {
	return &NamedEngine{
		name:   name,
		Engine: buildEngine(),
	}
}
