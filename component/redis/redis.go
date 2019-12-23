// redis 包提供redis客户端的连接与初始化
package redis

import (
	"github.com/ebar-go/ego/helper"
	"github.com/go-redis/redis"
	"net"
	"strconv"
	"time"
)

const (
	defaultPort        = 6379
	defaultPoolSize    = 100
	defaultMaxRetries  = 3
	defaultIdleTimeout = 10 * time.Second
)

// Conf redis配置
type Conf struct {
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

	// 是否为默认连接
	Default bool

	// 超时, 默认10s
	IdleTimeout time.Duration
}

// complete 自动补全
func (conf *Conf) complete() {
	conf.Port = helper.DefaultInt(conf.Port, defaultPort)
	conf.PoolSize = helper.DefaultInt(conf.PoolSize, defaultPoolSize)
	conf.MaxRetries = helper.DefaultInt(conf.MaxRetries, defaultMaxRetries)

	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = defaultIdleTimeout
	}
}

// Open 初始化连接池
func Open(conf Conf) (*redis.Client, error) {
	conf.complete()

	address := net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port))

	client := redis.NewClient(&redis.Options{
		Addr:        address,
		Password:    conf.Auth,
		PoolSize:    conf.PoolSize,    // Redis连接池大小
		MaxRetries:  conf.MaxRetries,  // 最大重试次数
		IdleTimeout: conf.IdleTimeout, // 空闲链接超时时间
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}
