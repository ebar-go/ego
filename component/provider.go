package component

import "sync"

const (
	componentCache           = "cache"
	componentLogger          = "logger"
	componentConfig          = "config"
	componentCurl            = "curl"
	componentJwt             = "jwt"
	componentTracer          = "tracer"
	componentEventDispatcher = "event-dispatcher"
	componentValidator       = "validator"
	componentRedis           = "redis"
	componentGorm            = "gorm"
)

// IProvider is a component provider
type IProvider interface {
	Logger() *Logger
	Cache() *Cache
	Config() *Config
	Curl() *Curl
	Jwt() *JWT
	Tracer() *Tracer
	EventDispatcher() *EventDispatcher
	Redis() *Redis
	Validator() *Validator
	Gorm() *Gorm
	Register(components ...Component)
	Get(name string) (Component, bool)
}

var providerInstance struct {
	once     sync.Once
	provider IProvider
}

// Initialize sets the component provider instance
func Initialize(provider IProvider) {
	providerInstance.once.Do(func() {
		providerInstance.provider = provider
	})
}

// Provider returns the  component provider singleton instance
func Provider() IProvider {
	if providerInstance.provider == nil {
		providerInstance.provider = NewContainer()
	}
	return providerInstance.provider
}

type Container struct {
	rmu        sync.RWMutex
	components map[string]Component
}

func (c *Container) build(name string) Component {
	switch name {
	case componentCache:
		return NewCache()
	case componentConfig:
		return NewConfig()
	case componentCurl:
		return NewCurl()
	case componentGorm:
		return NewGorm()
	case componentJwt:
		return NewJWT()
	case componentEventDispatcher:
		return NewEventDispatcher()
	case componentLogger:
		return NewLogger()
	case componentRedis:
		return NewRedis()
	case componentTracer:
		return NewTracer()
	case componentValidator:
		return NewValidator()
	}
	return nil
}

func (c *Container) GetOrInit(name string) Component {
	c.rmu.Lock()
	defer c.rmu.Unlock()
	component, ok := c.components[name]
	if ok {
		return component
	}
	component = c.build(name)
	c.components[name] = component
	return component

}
func (c *Container) Cache() *Cache {
	return c.GetOrInit(componentCache).(*Cache)
}

func (c *Container) Logger() *Logger {
	return c.GetOrInit(componentLogger).(*Logger)
}

func (c *Container) Config() *Config {
	return c.GetOrInit(componentConfig).(*Config)
}

func (c *Container) Curl() *Curl {
	return c.GetOrInit(componentCurl).(*Curl)
}

func (c *Container) Jwt() *JWT {
	return c.GetOrInit(componentJwt).(*JWT)
}

func (c *Container) Tracer() *Tracer {
	return c.GetOrInit(componentTracer).(*Tracer)
}

func (c *Container) EventDispatcher() *EventDispatcher {
	return c.GetOrInit(componentEventDispatcher).(*EventDispatcher)
}

func (c *Container) Redis() *Redis {
	return c.GetOrInit(componentRedis).(*Redis)
}

func (c *Container) Validator() *Validator {
	return c.GetOrInit(componentValidator).(*Validator)
}

func (c *Container) Gorm() *Gorm {
	return c.GetOrInit(componentGorm).(*Gorm)
}
func (c *Container) register(component Component) {
	c.components[component.Name()] = component
}
func (c *Container) Register(components ...Component) {
	c.rmu.Lock()
	for _, component := range components {
		c.register(component)
	}
	c.rmu.Unlock()
}

func (c *Container) Get(name string) (Component, bool) {
	c.rmu.RLock()
	defer c.rmu.RUnlock()
	item, ok := c.components[name]
	return item, ok
}

func NewContainer() *Container {
	return &Container{components: map[string]Component{}}
}
