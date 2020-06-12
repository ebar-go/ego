package config

import (
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"log"
	"time"
)

// Config 配置
type Config struct {
	viper.Viper
	// DI容器
	container *dig.Container
}

// New 实例
func New() *Config  {
	return &Config{container:dig.New()}
}

// 加载配置文件
func (conf *Config) LoadConfig(path string) error {
	conf.SetConfigFile(path)

	return conf.ReadInConfig()
}

// Server
func (conf *Config) Server() (options *server) {
	if err := conf.container.Invoke(func(o *server) {
		options = o
	}); err != nil {
		options = &server{
			Name:               conf.GetString(systemNameKey),
			Port:               conf.GetInt(httpPortKey),
			MaxResponseLogSize: conf.GetInt(maxResponseLogSizeKey),
			LogPath:            conf.GetString(logPathKey),
			JwtSignKey:         []byte(conf.GetString(jwtSignKey)),
			TraceHeader:        conf.GetString(traceHeaderKey),
			HttpRequestTimeOut: conf.GetInt(httpRequestTimeoutKey),
			Debug:              conf.GetBool(debugKey),
		}

		_ = conf.container.Provide(func() *server {
			return options
		})
	}
	return
}

// mysql
func (conf *Config) Mysql() map[string]mysql {
	var items map[string]mysql
	if err := viper.UnmarshalKey(mysqlKey, &items); err != nil {
		log.Println("Read Mysql Config:", err.Error())
		return nil
	}

	return items
}

// Redis
func (conf *Config) Redis() (options *redisConfig){
	if err := conf.container.Invoke(func(o *redisConfig) {
		options = o
	}); err != nil {
		options = &redisConfig{
			Host:        conf.GetString(redisHostKey),
			Port:        conf.GetInt(redisPortKey),
			Auth:        conf.GetString(redisPassKey),
			PoolSize:    conf.GetInt(redisPoolSizeKey),
			MaxRetries:  conf.GetInt(redisMaxRetriesKey),
			IdleTimeout: time.Duration(conf.GetInt(redisIdleTimeoutKey)) * time.Second,
		}

		_ = conf.container.Provide(func() *redisConfig {
			return options
		})
	}
	return
}