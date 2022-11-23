package component

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

// Cache provide cache instance builder function
type Cache struct {
	Named
	opts     cacheOptions
	once     sync.Once
	instance *cache.Cache
}

// WithDefaultExpiration sets default expiration for cache component
func (c *Cache) WithDefaultExpiration(defaultExpiration time.Duration) *Cache {
	c.opts.defaultExpiration = defaultExpiration
	return c
}

// WithCleanupInterval sets default cleanup interval for cache component
func (c *Cache) WithCleanupInterval(cleanupInterval time.Duration) *Cache {
	c.opts.cleanupInterval = cleanupInterval
	return c
}

// Default returns the default cache instance
func (c *Cache) Default() *cache.Cache {
	c.once.Do(func() {
		c.instance = c.Build()
	})
	return c.instance
}

type cacheOptions struct {
	defaultExpiration, cleanupInterval time.Duration
}

type cacheOption func(options *cacheOptions)

// Build returns a new cache instance
func (c *Cache) Build(opts ...cacheOption) *cache.Cache {
	defaultOptions := c.opts
	for _, setter := range opts {
		setter(&defaultOptions)
	}
	return cache.New(defaultOptions.defaultExpiration, defaultOptions.cleanupInterval)
}

func NewCache() *Cache {
	return &Cache{Named: componentCache, opts: cacheOptions{defaultExpiration: time.Minute * 5, cleanupInterval: time.Minute * 10}}
}
