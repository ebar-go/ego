package tracer

import (
	"fmt"
	"github.com/ebar-go/ego/utils/structure"
	"github.com/petermattis/goid"
	uuid "github.com/satori/go.uuid"
	"sync"
)

var (
	Set     = tracer().Set
	Get     = tracer().Get
	Release = tracer().Release
)

// Instance generate the uuid for per goroutine, use to mark user requests.
type Instance struct {
	collections *structure.ConcurrentMap[string, string]
}

// key use goroutine id to generate unique identifier
func (tracer *Instance) key() string {
	return fmt.Sprintf("g%d", goid.Get())
}

// Set sets the uuid for this goroutine
func (tracer *Instance) Set(id string) {
	tracer.collections.Set(tracer.key(), id)
}

// Get returns the uuid for this goroutine, it will generate a unique string if it doesn't exist'
func (tracer *Instance) Get() (id string) {
	id, ok := tracer.collections.Get(tracer.key())
	if ok {
		return id
	}
	id = uuid.NewV4().String()
	tracer.Set(id)
	return
}

// Release remove the uuid of this goroutine
func (tracer *Instance) Release() {
	tracer.collections.Del(tracer.key())
}

func New() *Instance {
	return &Instance{collections: structure.NewConcurrentMap[string, string]()}
}

var tracerSingleton = struct {
	once     sync.Once
	instance *Instance
}{}

func tracer() *Instance {
	tracerSingleton.once.Do(func() {
		tracerSingleton.instance = New()
	})
	return tracerSingleton.instance
}
