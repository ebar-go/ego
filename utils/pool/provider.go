package pool

import "sync"

// Provider represents a object provider
type Provider[Object any] interface {
	Acquire() Object
	Release(obj Object)
}

type SyncPoolProvider[Object any] struct {
	pool *sync.Pool
}

func (provider *SyncPoolProvider[Object]) Acquire() Object {
	return provider.pool.Get().(Object)
}

func (provider *SyncPoolProvider[Object]) Release(obj Object) {
	provider.pool.Put(obj)
}

func NewSyncPoolProvider[Object any](constructor func() interface{}) Provider[Object] {
	return &SyncPoolProvider[Object]{pool: &sync.Pool{New: constructor}}
}
