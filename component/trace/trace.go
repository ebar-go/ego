package trace

import (
	"github.com/ebar-go/ego/constant"
	"github.com/ebar-go/ego/utils/strings"
	"github.com/petermattis/goid"
	"sync"
)

var (
	traceIds = map[int64]string{}
	rwm      sync.RWMutex
)

func Id() string {
	return constant.TraceIdPrefix + strings.UUID()
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
