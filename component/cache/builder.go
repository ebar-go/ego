package cache

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type options struct {
	defaultExpiration, cleanupInterval time.Duration
}

type Option func(options *options)

func WithExpiration(d time.Duration) Option {
	return func(options *options) {
		options.defaultExpiration = d
	}
}
func WithCleanupInterval(d time.Duration) Option {
	return func(options *options) {
		options.cleanupInterval = d
	}
}

// Builder provide cache builder instance
type Builder struct {
	options  options
	once     sync.Once
	instance *cache.Cache
}

// Default returns the default cache instance
func (c *Builder) Default() *cache.Cache {
	c.once.Do(func() {
		c.instance = cache.New(c.options.defaultExpiration, c.options.cleanupInterval)
	})
	return c.instance
}

// Build returns a new cache instance
func (c *Builder) Build(opts ...Option) *cache.Cache {
	o := c.options
	for _, setter := range opts {
		setter(&o)
	}
	return cache.New(o.defaultExpiration, o.cleanupInterval)
}

func New() *Builder {
	return &Builder{options: options{defaultExpiration: time.Minute * 5, cleanupInterval: time.Minute * 10}}
}
