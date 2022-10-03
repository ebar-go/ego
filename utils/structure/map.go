package structure

import (
	"github.com/ebar-go/ego/errors"
	"sync"
)

// ConcurrentMap represents map container
type ConcurrentMap[T any] struct {
	mu    sync.RWMutex
	items map[string]T
}

// Get return item and exist
func (container *ConcurrentMap[T]) Get(key string) (item T, exist bool) {
	container.mu.RLock()
	defer container.mu.RUnlock()
	item, exist = container.items[key]
	return
}

// Set sets map item
func (container *ConcurrentMap[T]) Set(key string, value T) error {
	container.mu.Lock()
	container.items[key] = value
	container.mu.Unlock()
	return nil
}

// Find returns item
// if key is not exist, returns not found error
func (container *ConcurrentMap[T]) Find(key string) (item T, err error) {
	var exist bool
	item, exist = container.Get(key)
	if !exist {
		err = errors.NotFound("not found")
	}
	return
}

func NewConcurrentMap[T any]() *ConcurrentMap[T] {
	return &ConcurrentMap[T]{items: make(map[string]T, 0)}
}
