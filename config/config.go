package config

import (
	"github.com/ebar-go/ego/utils/number"
	"github.com/ebar-go/ego/utils/strings"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"time"
)

var Container = dig.New()

func init() {
	viper.AutomaticEnv()
}

func Server() (options *ServerOptions) {
	if err := Container.Invoke(func(o *ServerOptions) {
		options = o
	}); err != nil {
		options = &ServerOptions{
			Name:               strings.Default(viper.GetString("SYSTEM_NAME"), "app"),
			Port:               number.DefaultInt(viper.GetInt("HTTP_PORT"), 8080),
			MaxResponseLogSize: number.DefaultInt(viper.GetInt("MAX_RESPONSE_LOG_SIZE"), 1000),
			LogPath:            strings.Default(viper.GetString("LOG_PATH"), "/tmp"),
			JwtSignKey:         []byte(viper.GetString("JWT_KEY")),
			TraceHeader:        strings.Default(viper.GetString("TRACE_HEADER"), "gateway-trace"),
		}

		_ = Container.Provide(func() *ServerOptions {
			return options
		})
	}
	return
}

func Redis() (options *RedisOptions) {
	if err := Container.Invoke(func(o *RedisOptions) {
		options = o
	}); err != nil {
		options = &RedisOptions{
			AutoConnect: viper.GetBool("REDIS_AUTO_CONNECT"),
			Host:        strings.Default(viper.GetString("REDIS_HOST"), "127.0.0.1"),
			Port:        number.DefaultInt(viper.GetInt("REDIS_PORT"), 6379),
			Auth:        viper.GetString("REDIS_AUTH"),
			PoolSize:    number.DefaultInt(viper.GetInt("REDIS_POOL_SIZE"), 100),
			MaxRetries:  number.DefaultInt(viper.GetInt("REDIS_MAX_RETRIES"), 3),
			IdleTimeout: time.Duration(number.DefaultInt(viper.GetInt("REDIS_IDLE_TIMEOUT"), 5)) * time.Second,
		}

		_ = Container.Provide(func() *RedisOptions {
			return options
		})
	}
	return
}

func Mns() (options *MnsOptions) {
	if err := Container.Invoke(func(o *MnsOptions) {
		options = o
	}); err != nil {
		options = &MnsOptions{
			Url:             viper.GetString("MNS_ENDPOINT"),
			AccessKeyId:     viper.GetString("MNS_ACCESSID"),
			AccessKeySecret: viper.GetString("MNS_ACCESSKEY"),
		}

		_ = Container.Provide(func() *MnsOptions {
			return options
		})
	}
	return
}

func Mysql() (options *MysqlOptions) {
	if err := Container.Invoke(func(o *MysqlOptions) {
		options = o
	}); err != nil {
		options = &MysqlOptions{
			AutoConnect:        viper.GetBool("MYSQL_AUTO_CONNECT"),
			Name:               viper.GetString("MYSQL_DATABASE"),
			Host:               strings.Default(viper.GetString("MYSQL_MASTER_HOST"), "127.0.0.1"),
			Port:               number.DefaultInt(viper.GetInt("MYSQL_MASTER_PORT"), 3306),
			User:               viper.GetString("MYSQL_MASTER_USER"),
			Password:           viper.GetString("MYSQL_MASTER_PASS"),
			MaxOpenConnections: number.DefaultInt(viper.GetInt("MYSQL_MAX_OPEN_CONNECTIONS"), 10),
			MaxIdleConnections: number.DefaultInt(viper.GetInt("MYSQL_MAX_IDLE_CONNECTIONS"), 40),
		}

		_ = Container.Provide(func() *MysqlOptions {
			return options
		})
	}
	return
}
