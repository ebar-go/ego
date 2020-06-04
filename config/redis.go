package config

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"net"
	"strconv"
	"time"
)

// redisConfig redis配置
type redisConfig struct {
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

// Options get redis options
func (conf *redisConfig) Options() *redis.Options {
	address := net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port))

	return &redis.Options{
		Addr:        address,
		Password:    conf.Auth,
		PoolSize:    conf.PoolSize,    // Redis连接池大小
		MaxRetries:  conf.MaxRetries,  // 最大重试次数
		IdleTimeout: conf.IdleTimeout, // 空闲链接超时时间
	}
}

const (
	redisHostKey        = "redis.host"
	redisPortKey        = "redis.port"
	redisPassKey        = "redis.pass"
	redisPoolSizeKey    = "redis.pool_size"
	redisMaxRetriesKey  = "redis.max_retries"
	redisIdleTimeoutKey = "redis.idle_timeout"
)

func init() {
	viper.SetDefault(redisHostKey, "127.0.0.1")
	viper.SetDefault(redisPortKey, 6379)
	viper.SetDefault(redisPoolSizeKey, 100)
	viper.SetDefault(redisMaxRetriesKey, 3)
	viper.SetDefault(redisIdleTimeoutKey, 5)
}

func Redis() (options *redisConfig) {
	if err := Container.Invoke(func(o *redisConfig) {
		options = o
	}); err != nil {
		options = &redisConfig{
			Host:        viper.GetString(redisHostKey),
			Port:        viper.GetInt(redisPortKey),
			Auth:        viper.GetString(redisPassKey),
			PoolSize:    viper.GetInt(redisPoolSizeKey),
			MaxRetries:  viper.GetInt(redisMaxRetriesKey),
			IdleTimeout: time.Duration(viper.GetInt(redisIdleTimeoutKey)) * time.Second,
		}

		_ = Container.Provide(func() *redisConfig {
			return options
		})
	}
	return
}
