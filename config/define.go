package config

import (
	"fmt"
	"github.com/go-redis/redis"
	"net"
	"strconv"
	"time"
)

// ServerOptions 服务配置
type ServerOptions struct {
	// 服务名称
	Name string

	// 服务端口号
	Port int

	// 响应日志最大长度
	MaxResponseLogSize int

	// 日志路径
	LogPath string
	// jwt的key
	JwtSignKey []byte

	// trace header key
	TraceHeader string

	// http request timeout
	HttpRequestTimeOut int
}

// MnsOptions 阿里云MNS 配置项
type MnsOptions struct {
	Url             string `json:"url"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

// Dsn return mysql dsn
func (options MysqlOptions) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		options.User,
		options.Password,
		net.JoinHostPort(options.Host, strconv.Itoa(options.Port)),
		options.Name)
}

// MysqlOptions
type MysqlOptions struct {
	// auto connect
	AutoConnect bool

	// database name
	Name string

	// host
	Host string

	// port, default 3306
	Port int

	// user, default root
	User string

	// password
	Password string

	// log mode
	LogMode bool

	// MaxIdleConnections, default 10
	MaxIdleConnections int

	// MaxOpenConnections, default 40
	MaxOpenConnections int

	// max life time, default 10
	MaxLifeTime int
}

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