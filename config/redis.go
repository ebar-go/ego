package config

import (
	"github.com/ebar-go/ego/utils/number"
	"github.com/go-redis/redis"
	"net"
	"strconv"
	"time"
)

// RedisOptions redis配置
type RedisOptions struct {
	// AutoConnect
	AutoConnect bool

	// host
	Host string

	// port, default 6379
	Port int

	// auth
	Auth string

	// pool size, default 100
	PoolSize int

	// max retries, default 3
	MaxRetries int

	// timeout, default 10 seconds
	IdleTimeout time.Duration
}

const (
	redisDefaultPort        = 6379
	redisDefaultPoolSize    = 100
	redisDefaultMaxRetries  = 3
	redisDefaultIdleTimeout = 10 * time.Second
)

// complete set default config
func (conf *RedisOptions) complete() {
	conf.Port = number.DefaultInt(conf.Port, redisDefaultPort)
	conf.PoolSize = number.DefaultInt(conf.PoolSize, redisDefaultPoolSize)
	conf.MaxRetries = number.DefaultInt(conf.MaxRetries, redisDefaultMaxRetries)

	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = redisDefaultIdleTimeout
	}
}

// Options get redis options
func (conf *RedisOptions) Options() *redis.Options {
	address := net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port))

	return &redis.Options{
		Addr:        address,
		Password:    conf.Auth,
		PoolSize:    conf.PoolSize,    // Redis连接池大小
		MaxRetries:  conf.MaxRetries,  // 最大重试次数
		IdleTimeout: conf.IdleTimeout, // 空闲链接超时时间
	}
}
