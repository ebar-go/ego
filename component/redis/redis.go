// redis 包提供redis客户端的连接与初始化
package redis

import (
	"github.com/go-redis/redis"
	"time"
	"sync"
	"net"
	"strconv"
	"fmt"
	"os"
)


var redisConnection *redis.Client
var redisInitOnce sync.Once

const (
	defaultPort = 6379
	defaultPoolSize = 100
	defaultMaxRetries = 3
	defaultIdleTimeout = 10*time.Second
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

	// 超时, 默认10s
	IdleTimeout time.Duration
}

// complete 自动补全
func (conf *Conf) complete() {
	if conf.Port == 0 {
		conf.Port = defaultPort
	}

	if conf.PoolSize == 0 {
		conf.PoolSize = defaultPoolSize
	}

	if conf.MaxRetries == 0 {
		conf.MaxRetries = defaultMaxRetries
	}

	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = defaultIdleTimeout
	}
}

// InitPool 初始化连接池
func InitPool(conf Conf)  (err error) {
	redisInitOnce.Do(func() {
		conf.complete()

		address := net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port))

		client := redis.NewClient(&redis.Options{
			Addr: address,
			Password:conf.Auth,
			PoolSize: conf.PoolSize,  // Redis连接池大小
			MaxRetries: conf.MaxRetries,              // 最大重试次数
			IdleTimeout: conf.IdleTimeout,            // 空闲链接超时时间
		})

		if _, err := client.Ping().Result();err != nil {
			fmt.Println("Redis连接失败:", err)
			os.Exit(-1)
		}

		redisConnection = client
		fmt.Println("连接Redis成功:",address)
	})

	return nil
}

// GetConnection 获取连接
func GetConnection() *redis.Client {
	return redisConnection
}


