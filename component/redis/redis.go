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


const (
	defaultPort = 6379
	defaultPoolSize = 100
	defaultMaxRetries = 3
	defaultIdleTimeout = 10*time.Second
)

//
var group *ConnectionGroup

func init() {
	group = &ConnectionGroup{
		lock:        &sync.Mutex{},
		connections: make(map[string]*redis.Client),
	}
}

// GetConnectionGroup 获取连接池组
func GetConnectionGroup() *ConnectionGroup {
	return group
}

// ConnectionGroup 数据库连接组
type ConnectionGroup struct {
	lock        *sync.Mutex
	defaultKey  string
	connections map[string]*redis.Client
}

type ConnectFailedHandler func(err error)

// Conf redis配置
type Conf struct {
	Key string
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

	ConnectFailedHandler ConnectFailedHandler

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

	if conf.ConnectFailedHandler == nil {
		conf.ConnectFailedHandler = func(err error) {
			fmt.Printf("redis %s connect failed:%s", conf.Key, err.Error())
			os.Exit(-1)
		}
	}
}

// InitPool 初始化连接池
func InitPool(confItems ...Conf)  (err error) {
	// 加锁
	group.lock.Lock()

	defaultConnectionKey := ""
	for key, conf := range confItems {
		// 如果没有设置default选项，则默认取第一个
		if key == 0 {
			defaultConnectionKey = conf.Key
		}

		if conf.Default {
			defaultConnectionKey = conf.Key
		}

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
			conf.ConnectFailedHandler(err)
		}

		group.connections[conf.Key] = client
		fmt.Println("连接Redis成功:",address)
	}

	group.defaultKey = defaultConnectionKey

	return nil
}

// GetConnection 获取默认连接
func GetConnection() *redis.Client {
	return GetConnectionByKey(group.defaultKey)
}

// GetConnection 获取连接
func GetConnectionByKey(key string) *redis.Client {
	return group.connections[key]
}


