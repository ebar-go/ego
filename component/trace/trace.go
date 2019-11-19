package trace

import (
	"sync"
	"github.com/petermattis/goid"
)

var (
	traceIds = map[int64]string{}
	rwm        sync.RWMutex
)

// SetTraceId
func SetTraceId(id string) {
	goID := getGoroutineId()
	rwm.Lock()
	defer rwm.Unlock()

	traceIds[goID] = id
}

// GetTraceId
func GetTraceId() string {
	goID := getGoroutineId()
	rwm.RLock()
	defer rwm.RUnlock()

	return traceIds[goID]
}

// DeleteTraceId
func DeleteTraceId() {
	goID := getGoroutineId()
	rwm.Lock()
	defer rwm.Unlock()

	delete(traceIds, goID)
}

func getGoroutineId() int64 {
	return goid.Get()
}
