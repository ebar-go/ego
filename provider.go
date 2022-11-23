package ego

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/component/cache"
	"sync"
)

type Provider struct {
	cc component.Instance[*cache.Builder]
}

var providerInstance = struct {
	once     sync.Once
	instance *Provider
}{}

func provider() *Provider {
	providerInstance.once.Do(func() {
		providerInstance.instance = &Provider{
			cc: component.NewInstance[*cache.Builder]("default", cache.New()),
		}
	})
	return providerInstance.instance
}

func CacheBuilder() *cache.Builder {
	return provider().cc.Delegate()
}
