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

func Initialize(provider ComponentProvider) {
	providerInstance.once.Do(func() {
		providerInstance.provider = provider
	})
}
func Provider() ComponentProvider {
	if providerInstance.provider == nil {
		providerInstance.provider = NewContainer()
	}
	return providerInstance.provider
}

type Container struct {
	cache  *Cache
	logger *Logger
}

func (c *Container) Cache() *Cache {
	if c.cache == nil {
		c.cache = NewCache()
	}
	return c.cache
}

func (c *Container) Logger() *Logger {
	if c.logger == nil {
		c.logger = NewLogger()
	}
	return c.logger
}

func (c *Container) register(component Component) {
	if cache, ok := component.(*Cache); ok {
		c.cache = cache
	}
}
func (c *Container) Register(components ...Component) {
	for _, component := range components {
		c.register(component)
	}
}

func NewContainer() *Container {
	return &Container{}
}
