package component

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

// Cache provide cache instance builder function
type Cache struct {
	Named
	defaultExpiration, cleanupInterval time.Duration
	once                               sync.Once
	instance                           *cache.Cache
}

// WithDefaultExpiration sets default expiration for cache component
func (c *Cache) WithDefaultExpiration(defaultExpiration time.Duration) *Cache {
	c.defaultExpiration = defaultExpiration
	return c
}

// WithCleanupInterval sets default cleanup interval for cache component
func (c *Cache) WithCleanupInterval(cleanupInterval time.Duration) *Cache {
	c.cleanupInterval = cleanupInterval
	return c
}

// Default returns the default cache instance
func (c *Cache) Default() *cache.Cache {
	c.once.Do(func() {
		c.instance = c.Build()
	})
	return c.instance
}

// Build returns a new cache instance
func (c *Cache) Build() *cache.Cache {
	return cache.New(c.defaultExpiration, c.cleanupInterval)
}

// BuildWith returns a new cache instance with the given expiration.
func (c *Cache) BuildWith(defaultExpiration, cleanupInterval time.Duration) *cache.Cache {
	return cache.New(defaultExpiration, cleanupInterval)
}

func NewCache() *Cache {
	return &Cache{Named: componentCache, defaultExpiration: time.Minute * 5, cleanupInterval: time.Minute * 10}
}
