package component

import "sync"

// ComponentProvider is a component provider
type ComponentProvider interface {
	Logger() *Logger
	Cache() *Cache
	Config() *Config
	Get(name string) (Component, bool)
}

var providerInstance struct {
	once     sync.Once
	provider ComponentProvider
}

// Initialize sets the component provider instance
func Initialize(provider ComponentProvider) {
	providerInstance.once.Do(func() {
		providerInstance.provider = provider
	})
}

// Provider returns the  component provider singleton instance
func Provider() ComponentProvider {
	if providerInstance.provider == nil {
		providerInstance.provider = NewContainer()
	}
	return providerInstance.provider
}

type Container struct {
	cache  *Cache
	logger *Logger
	config *Config
	rmu    sync.RWMutex
	others map[string]Component
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

func (c *Container) Config() *Config {
	if c.config == nil {
		c.config = NewConfig()
	}
	return c.config
}

func (c *Container) register(component Component) {
	if cache, ok := component.(*Cache); ok {
		c.cache = cache
		return
	}

	if config, ok := component.(*Config); ok {
		c.config = config
		return
	}

	if logger, ok := component.(*Logger); ok {
		c.logger = logger
		return
	}

	c.rmu.Lock()
	c.others[component.Name()] = component
	c.rmu.Unlock()

}
func (c *Container) Register(components ...Component) {
	for _, component := range components {
		c.register(component)
	}
}

func (c *Container) Get(name string) (Component, bool) {
	c.rmu.RLock()
	defer c.rmu.RUnlock()
	item, ok := c.others[name]
	return item, ok
}

func NewContainer() *Container {
	return &Container{others: map[string]Component{}}
}
