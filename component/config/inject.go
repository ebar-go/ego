/**
 * @Author: Hongker
 * @Description:
 * @File:  inject
 * @Version: 1.0.0
 * @Date: 2021/4/3 18:00
 */

package config

import (
	"github.com/ebar-go/ego/component/etcd"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/mysql"
	"github.com/ebar-go/ego/component/redis"
	"github.com/ebar-go/ego/http"
	"go.uber.org/dig"
	"time"
)

func Inject(container *dig.Container) {
	_ = container.Provide(New)
	_ = container.Provide(httpConfig)
	_ = container.Provide(logConfig)
	_ = container.Provide(mysqlConfig)
	_ = container.Provide(redisConfig)
	_ = container.Provide(etcdConfig)
}
func httpConfig(conf *Config) *http.Config {
	return &http.Config{
		Port:        conf.GetDefaultInt("http.port", 8080),
		JwtSignKey:  []byte(conf.GetString("http.jwtSign")),
		TraceHeader: conf.GetDefaultString("http.traceHeader", "request-trace"),
		Pprof:       conf.GetBool("http.pprof"),
		Swagger:     conf.GetBool("http.swagger"),
	}
}

func logConfig(conf *Config) *log.Config {
	return &log.Config{
		Path:  conf.GetDefaultString("log.path", "path"),
		Debug: conf.GetBool("log.debug"),
	}
}

func mysqlConfig(conf *Config) *mysql.Config {
	return &mysql.Config{
		MaxIdleConnections: conf.GetInt("mysql.maxIdleConnections"),
		MaxOpenConnections: conf.GetInt("mysql.maxOpenConnections"),
		MaxLifeTime:        conf.GetInt("mysql.maxLifeTime"),
		Dsn:                conf.GetString("mysql.dsn"),
	}
}

func redisConfig(conf *Config) *redis.Config {
	return &redis.Config{
		Host:        conf.GetDefaultString("redis.host", "127.0.0.1"),
		Port:        conf.GetDefaultInt("redis.port", 6379),
		Auth:        conf.GetString("redis.pass"),
		PoolSize:    conf.GetDefaultInt("redis.poolSize", 100),
		MaxRetries:  conf.GetDefaultInt("redis.maxRetries", 3),
		IdleTimeout: time.Duration(conf.GetDefaultInt("redis.idleTimeout", 10)) * time.Second,
		Cluster:     conf.GetStringSlice("redis.cluster"),
	}
}

func etcdConfig(conf *Config) *etcd.Config {
	return &etcd.Config{
		Endpoints: conf.GetStringSlice("etcd.endpoints"),
		Timeout:   time.Second * time.Duration(conf.GetDefaultInt("etcd.timeout", 10)),
	}
}
