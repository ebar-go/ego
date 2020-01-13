package config

import (
	"github.com/ebar-go/ego/utils/conv"
	"github.com/ebar-go/ego/utils/number"
	"github.com/ebar-go/ego/utils/strings"
)

// Config接口
type Config interface {
	Server() *ServerOption
	Redis() *RedisOptions
	Mysql() *MysqlOptions
	Mns() *MnsConfig
}

// config 系统配置项
type config struct {
	serverOption *ServerOption

	// redis config
	redisConfig *RedisOptions

	// mysql config
	mysqlOptions *MysqlOptions

	// mns config
	mnsConfig *MnsConfig
}



func (conf *config) Server() *ServerOption {
	return conf.serverOption
}

// Redis config
func (config *config) Redis() *RedisOptions {
	return config.redisConfig
}

// Mysql config
func (config *config) Mysql() *MysqlOptions {
	return config.mysqlOptions
}

// Mns config
func (config *config) Mns() *MnsConfig {
	return config.mnsConfig
}

// LoadEnv 通过加载环境变量初始化系统配置
func LoadEnv() Config {
	instance := &config{
		serverOption: new(ServerOption),
		redisConfig:  new(RedisOptions),
		mysqlOptions: new(MysqlOptions),
		mnsConfig:    new(MnsConfig),
	}
	instance.serverOption.Name = strings.Default(Getenv("SYSTEM_NAME"), "app")
	instance.serverOption.Port = number.DefaultInt(conv.String2Int(Getenv("HTTP_PORT")), 8080)

	instance.serverOption.LogPath = strings.Default(Getenv("LOG_PATH"), "/tmp")
	instance.serverOption.MaxResponseLogSize = number.DefaultInt(conv.String2Int(Getenv("MAX_RESPONSE_LOG_SIZE")), 1000)

	instance.serverOption.JwtSignKey = []byte(Getenv("JWT_KEY"))
	instance.serverOption.TraceHeader = strings.Default(Getenv("TRACE_HEADER"), "gateway-trace")

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
