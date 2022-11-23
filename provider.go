package ego

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/component/cache"
	"github.com/ebar-go/ego/component/config"
	"sync"
)

type Provider struct {
	cb component.Instance[*cache.Builder]
	cc component.Instance[*config.Config]
}

var providerInstance = struct {
	once     sync.Once
	instance *Provider
}{}

func provider() *Provider {
	providerInstance.once.Do(func() {
		providerInstance.instance = &Provider{
			cb: component.NewInstance[*cache.Builder]("default", cache.New()),
			cc: component.NewInstance[*config.Config]("default", config.New()),
		}
	})
	return providerInstance.instance
}

func CacheBuilder() *cache.Builder {
	return provider().cb.Delegate()
}

func Config() *config.Config {
	return provider().cc.Delegate()
}
