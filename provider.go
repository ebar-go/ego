package ego

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/component/cache"
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/curl"
	"github.com/robfig/cron"
	"sync"
)

type Provider struct {
	cb   component.Instance[*cache.Builder]
	cc   component.Instance[*config.Config]
	cron component.Instance[*cron.Cron]
	curl component.Instance[*curl.Curl]
}

var providerInstance = struct {
	once     sync.Once
	instance *Provider
}{}

func provider() *Provider {
	providerInstance.once.Do(func() {
		name := "default"
		providerInstance.instance = &Provider{
			cb:   component.NewInstance[*cache.Builder](name, cache.New()),
			cc:   component.NewInstance[*config.Config](name, config.New()),
			cron: component.NewInstance[*cron.Cron](name, cron.New()),
			curl: component.NewInstance[*curl.Curl](name, curl.New()),
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

func Cron() *cron.Cron {
	return provider().cron.Delegate()
}

func Curl() *curl.Curl {
	return provider().curl.Delegate()
}