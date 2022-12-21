package structure

import (
	"github.com/ebar-go/ego/errors"
	"sync"
	"sync/atomic"
)

type KeyType interface {
	int | int16 | int32 | int64 | string
}

type Map[Key KeyType, Val any] interface {
	Get(key Key) (item Val, exist bool)
	Set(key Key, value Val)
	Del(key Key)
	Iterator(fn func(key Key, val Val))
	Len() int
}

// ConcurrentMap represents safe concurrent map container
type ConcurrentMap[Key KeyType, Val any] struct {
	mu    sync.RWMutex
	items map[Key]Val
}

// Get return item and exist
func (container *ConcurrentMap[Key, T]) Get(key Key) (item T, exist bool) {
	container.mu.RLock()
	defer container.mu.RUnlock()
	item, exist = container.items[key]
	return
}

// Set sets map item
func (container *ConcurrentMap[Key, T]) Set(key Key, value T) {
	container.mu.Lock()
	container.items[key] = value
	container.mu.Unlock()
}

// Find returns item
// if key is not exist, returns not found error
func (container *ConcurrentMap[Key, T]) Find(key Key) (item T, err error) {
	var exist bool
	item, exist = container.Get(key)
	if !exist {
		err = errors.NotFound("not found")
	}
	return
}

// Del remove keys from container
func (container *ConcurrentMap[Key, T]) Del(key Key) {
	container.mu.Lock()
	delete(container.items, key)
	container.mu.Unlock()
}

func (container *ConcurrentMap[Key, T]) Iterator(fn func(key Key, val T)) {
	container.mu.RLock()
	for key, val := range container.items {
		fn(key, val)
	}
	container.mu.RUnlock()
}

func (container *ConcurrentMap[Key, T]) Len() int {
	container.mu.RLock()
	defer container.mu.RUnlock()
	return len(container.items)
}
func NewConcurrentMap[Key KeyType, T any]() *ConcurrentMap[Key, T] {
	return &ConcurrentMap[Key, T]{items: make(map[Key]T, 0)}
}

type LockFreeMap[Key KeyType, Val any] struct {
	lock  atomic.Value
	items map[Key]Val
}

func NewLockFreeMap[Key KeyType, T any]() *LockFreeMap[Key, T] {
	return &LockFreeMap[Key, T]{items: make(map[Key]T, 0)}
}

// Get return item and exist
func (container *LockFreeMap[Key, T]) Get(key Key) (item T, exist bool) {
	container.lock.Store(true)
	item, exist = container.items[key]
	return
}

// Set sets map item
func (container *LockFreeMap[Key, T]) Set(key Key, value T) {
	container.lock.Store(true)
	container.items[key] = value
}

// Find returns item
// if key is not exist, returns not found error
func (container *LockFreeMap[Key, T]) Find(key Key) (item T, err error) {
	var exist bool
	item, exist = container.Get(key)
	if !exist {
		err = errors.NotFound("not found")
	}
	return
}

// Del remove keys from container
func (container *LockFreeMap[Key, T]) Del(key Key) {
	container.lock.Store(true)
	delete(container.items, key)
}

func (container *LockFreeMap[Key, T]) Iterator(fn func(key Key, val T)) {
	container.lock.Store(true)
	for key, val := range container.items {
		fn(key, val)
	}
}

func (container *LockFreeMap[Key, T]) Len() int {
	container.lock.Store(true)
	return len(container.items)
}
