package redis

import (
	"github.com/go-redis/redis"
	"net"
	"strconv"
	"time"
)

// Config Redis配置
type Config struct {
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

	// 集群
	Cluster []string
}

// Options 单机选项
func (conf *Config) Options() *redis.Options {
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
func (conf *Config) ClusterOption() *redis.ClusterOptions {
	return &redis.ClusterOptions{
		Addrs:       conf.Cluster,
		Password:    conf.Auth,
		PoolSize:    conf.PoolSize,    // Redis连接池大小
		MaxRetries:  conf.MaxRetries,  // 最大重试次数
		IdleTimeout: conf.IdleTimeout, // 空闲链接超时时间
	}
}
