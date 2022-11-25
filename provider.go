package ego

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/component/cache"
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/curl"
	"github.com/ebar-go/ego/component/event"
	"github.com/ebar-go/ego/component/logger"
	"github.com/ebar-go/ego/component/redis"
	"github.com/robfig/cron"
	"sync"
)

type Provider struct {
	cb     component.Instance[*cache.Builder]
	cc     component.Instance[*config.Config]
	cron   component.Instance[*cron.Cron]
	curl   component.Instance[*curl.Curl]
	event  component.Instance[*event.Dispatcher]
	redis  component.Instance[*redis.Instance]
	logger component.Instance[*logger.Logger]
}

var providerInstance = struct {
	once     sync.Once
	instance *Provider
}{}

func provider() *Provider {
	providerInstance.once.Do(func() {
		name := "default"
		providerInstance.instance = &Provider{
			cb:     component.NewInstance[*cache.Builder](name, cache.New()),
			cc:     component.NewInstance[*config.Config](name, config.New()),
			cron:   component.NewInstance[*cron.Cron](name, cron.New()),
			curl:   component.NewInstance[*curl.Curl](name, curl.New()),
			event:  component.NewInstance[*event.Dispatcher](name, event.New()),
			redis:  component.NewInstance[*redis.Instance](name, redis.New()),
			logger: component.NewInstance[*logger.Logger](name, logger.New()),
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

func Event() *event.Dispatcher {
	return provider().event.Delegate()
}

// ListenEvent listen with type parameters
func ListenEvent[T any](eventName string, handler func(param T)) {
	Event().Listen(eventName, func(param any) {
		data, ok := param.(T)
		if !ok {
			return
		}
		handler(data)
	})
}

func Redis() *redis.Instance {
	return provider().redis.Delegate()
}

func Logger() *logger.Logger {
	return provider().logger.Delegate()
}
