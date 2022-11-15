package component

import (
	"fmt"
	"github.com/ebar-go/ego/utils/structure"
	"github.com/petermattis/goid"
	uuid "github.com/satori/go.uuid"
)

// Tracer generate the uuid for per goroutine, use to mark user requests.
type Tracer struct {
	Named
	collections *structure.ConcurrentMap[string]
}

// key use goroutine id to generate unique identifier
func (tracer *Tracer) key() string {
	return fmt.Sprintf("g%d", goid.Get())
}

// Set sets the uuid for this goroutine
func (tracer *Tracer) Set(id string) {
	tracer.collections.Set(tracer.key(), id)
}

// Get returns the uuid for this goroutine, it will generate a unique string if it doesn't exist'
func (tracer *Tracer) Get() (id string) {
	id, ok := tracer.collections.Get(tracer.key())
	if ok {
		return id
	}
	id = uuid.NewV4().String()
	tracer.Set(id)
	return
}

// Release remove the uuid of this goroutine
func (tracer *Tracer) Release() {
	tracer.collections.Del(tracer.key())
}

func NewTracer() *Tracer {
	return &Tracer{Named: componentTracer, collections: structure.NewConcurrentMap[string]()}
}
