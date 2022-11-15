package cmap

import "sync"

// Container represents a safe map
type Container[Key int | int16 | string, Value any] struct {
	rmu   sync.RWMutex
	items map[Key]Value
}

// Set sets the value of a key in the container
func (c *Container[Key, Value]) Set(key Key, value Value) {
	c.rmu.Lock()
	c.items[key] = value
	c.rmu.Unlock()
}

// Get returns the value associated
func (c *Container[Key, Value]) Get(key Key) (item Value, exist bool) {
	c.rmu.RLock()
	item, exist = c.items[key]
	c.rmu.RUnlock()
	return
}

// Del removes the item from the container
func (c *Container[Key, Value]) Del(key Key) {
	c.rmu.Lock()
	delete(c.items, key)
	c.rmu.Unlock()
}

// NewContainer creates a new Container instance
func NewContainer[Key int | int16 | string, Value any]() *Container[Key, Value] {
	return &Container[Key, Value]{
		items: make(map[Key]Value),
	}
}

type ShardContainer[T any] []T

func NewShardContainer[T any](size int, constructor func() T) ShardContainer[T] {
	sc := make(ShardContainer[T], size)
	for i := 0; i < size; i++ {
		sc[i] = constructor()
	}
	return sc
}

func (shard ShardContainer[T]) GetShard(num int) T {
	idx := num & (len(shard) - 1)
	return shard[idx]
}

func (shard ShardContainer[T]) Iterator(iterator func(T)) {
	for _, item := range shard {
		iterator(item)
	}
}
