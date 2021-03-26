package trace

import (
	"github.com/ebar-go/egu"
	"github.com/petermattis/goid"
	"sync"
)

var (
	traceIds = map[int64]string{}
	rwm      sync.RWMutex
)

const (
	prefix = "trace:"
)

func Id() string {
	return prefix + egu.UUID()
}

// SetTraceId
func Set(id string) {
	goID := getGoroutineId()
	rwm.Lock()
	defer rwm.Unlock()

	traceIds[goID] = id
}

// Get
func Get() string {
	goID := getGoroutineId()
	rwm.RLock()
	defer rwm.RUnlock()

	return traceIds[goID]
}

// GC
func GC() {
	goID := getGoroutineId()
	rwm.Lock()
	defer rwm.Unlock()

	delete(traceIds, goID)
}

func getGoroutineId() int64 {
	return goid.Get()
}

func Go(f func())  {
	go func(traceId string) {
		Set(traceId)
		defer GC()
		f()
	}(Get())
}
