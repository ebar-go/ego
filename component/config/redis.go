package config

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"net"
	"strconv"
	"strings"
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

	Cluster string
}

// Options 单机选项
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

// ClusterOption 集群选项
func (conf *redisConfig) ClusterOption() *redis.ClusterOptions {
	return &redis.ClusterOptions{
		Addrs: strings.Split(conf.Cluster, ","),
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
	redisCluster = "redis.cluster"
)

func init() {
	viper.SetDefault(redisHostKey, "127.0.0.1")
	viper.SetDefault(redisPortKey, 6379)
	viper.SetDefault(redisPoolSizeKey, 100)
	viper.SetDefault(redisMaxRetriesKey, 3)
	viper.SetDefault(redisIdleTimeoutKey, 5)
}
