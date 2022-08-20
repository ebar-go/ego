package component

import "sync"

type ComponentProvider interface {
	Logger() *Logger
	Cache() *Cache
}

var providerInstance struct {
	once     sync.Once
	provider ComponentProvider
}

func Provider() ComponentProvider {
	providerInstance.once.Do(func() {

	})
	return providerInstance.provider
}
