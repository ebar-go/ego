package config

import (
	"github.com/ebar-go/ego/utils/number"
	"github.com/go-redis/redis"
	"net"
	"strconv"
	"time"
)

// RedisConfig redis配置
type RedisConfig struct {
	// 地址
	Host string

	// 端口号
	Port int

	// 密码
	Auth string

	// 连接池大小,默认100个连接
	PoolSize int

	// 最大尝试次数,默认3次
	MaxRetries int

	// 超时, 默认10s
	IdleTimeout time.Duration
}

const (
	redisDefaultPort        = 6379
	redisDefaultPoolSize    = 100
	redisDefaultMaxRetries  = 3
	redisDefaultIdleTimeout = 10 * time.Second
)

// complete set default config
func (conf *RedisConfig) complete() {
	conf.Port = number.DefaultInt(conf.Port, redisDefaultPort)
	conf.PoolSize = number.DefaultInt(conf.PoolSize, redisDefaultPoolSize)
	conf.MaxRetries = number.DefaultInt(conf.MaxRetries, redisDefaultMaxRetries)

	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = redisDefaultIdleTimeout
	}
}

// Options get redis options
func (conf *RedisConfig) Options() *redis.Options {
	address := net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port))

	return &redis.Options{
		Addr:        address,
		Password:    conf.Auth,
		PoolSize:    conf.PoolSize,    // Redis连接池大小
		MaxRetries:  conf.MaxRetries,  // 最大重试次数
		IdleTimeout: conf.IdleTimeout, // 空闲链接超时时间
	}
}
