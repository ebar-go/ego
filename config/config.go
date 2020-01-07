package config

import (
	"github.com/ebar-go/ego/utils/conv"
	"github.com/ebar-go/ego/utils/number"
	"github.com/ebar-go/ego/utils/strings"
)

// Config 系统配置项
type Config struct {
	// 服务名称
	ServiceName string

	// 服务端口号
	ServicePort int

	// 响应日志最大长度
	MaxResponseLogSize int

	// 日志路径
	LogPath string

	// jwt的key
	JwtSignKey []byte

	// trace header key
	TraceHeader string

	// redis config
	redisConfig *RedisOptions

	// mysql config
	mysqlOptions *MysqlOptions

	// mns config
	mnsConfig *MnsConfig
}

// Redis config
func (config *Config) Redis() *RedisOptions {
	return config.redisConfig
}

// Mysql config
func (config *Config) Mysql() *MysqlOptions {
	return config.mysqlOptions
}

// Mns config
func (config *Config) Mns() *MnsConfig {
	return config.mnsConfig
}

// init 通过读取环境变量初始化系统配置
func NewInstance() *Config {
	instance := &Config{}
	instance.ServiceName = strings.Default(Getenv("SYSTEM_NAME"), "app")
	instance.ServicePort = number.DefaultInt(conv.String2Int(Getenv("HTTP_PORT")), 8080)

	instance.LogPath = strings.Default(Getenv("LOG_PATH"), "/tmp")
	instance.MaxResponseLogSize = number.DefaultInt(conv.String2Int(Getenv("MAX_RESPONSE_LOG_SIZE")), 1000)

	instance.JwtSignKey = []byte(Getenv("JWT_KEY"))
	instance.TraceHeader = strings.Default(Getenv("TRACE_HEADER"), "gateway-trace")

	// init mysql config
	instance.redisConfig = &RedisOptions{
		AutoConnect:strings.ToBool(Getenv("REDIS_AUTO_CONNECT")),
		Host: strings.Default(Getenv("REDIS_HOST"), "127.0.0.1"),
		Port: number.DefaultInt(conv.String2Int(Getenv("REDIS_PORT")), 6379),
		Auth: Getenv("REDIS_AUTH"),
	}
	instance.redisConfig.complete()

	// init redis config
	instance.mysqlOptions = &MysqlOptions{
		AutoConnect: strings.ToBool(Getenv("MYSQL_AUTO_CONNECT")),
		Name:     Getenv("MYSQL_DATABASE"),
		Host:     strings.Default(Getenv("MYSQL_MASTER_HOST"), "127.0.0.1"),
		Port:     number.DefaultInt(conv.String2Int(Getenv("MYSQL_MASTER_PORT")), 3306),
		User:     Getenv("MYSQL_MASTER_USER"),
		Password: Getenv("MYSQL_MASTER_PASS"),
	}
	instance.mysqlOptions.complete()

	// mns config
	instance.mnsConfig = &MnsConfig{
		Url:             Getenv("MNS_ENDPOINT"),
		AccessKeyId:     Getenv("MNS_ACCESSID"),
		AccessKeySecret: Getenv("MNS_ACCESSKEY"),
	}

	return instance
}
