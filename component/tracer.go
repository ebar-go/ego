package component

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/petermattis/goid"
	uuid "github.com/satori/go.uuid"
)

// Tracer generate the uuid for per goroutine
type Tracer struct {
	Named
	collections cmap.ConcurrentMap
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
func (tracer *Tracer) Get() string {
	val, ok := tracer.collections.Get(tracer.key())
	if ok {
		return val.(string)
	}
	id := uuid.NewV4().String()
	tracer.Set(id)
	return id
}

// Release remove the uuid of this goroutine
func (tracer *Tracer) Release() {
	tracer.collections.Remove(tracer.key())
}

func NewTracer() *Tracer {
	return &Tracer{Named: "tracer", collections: cmap.New()}
}
