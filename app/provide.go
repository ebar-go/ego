package app

import (
	"github.com/ebar-go/ego/component/auth"
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/etcd"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/mysql"
	"github.com/ebar-go/ego/component/redis"
	"net"
	"net/http"
	"time"
)

// newHttpClient
func newHttpClient(conf *config.Config) *http.Client {
	return &http.Client{
		Transport: &http.Transport{ // 配置连接池
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			IdleConnTimeout: time.Duration(conf.Server().HttpRequestTimeOut) * time.Second,
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Duration(conf.Server().HttpRequestTimeOut) * time.Second,
	}
}

// newLogger
func newLogger (conf *config.Config) *log.Logger {
	return log.New(conf.Server().LogPath,
		conf.Server().Debug,
		map[string]interface{}{
			"system_name": conf.Server().Name,
		})
}

// newJwt
func newJwt(conf *config.Config) *auth.JwtAuth {
	return auth.New(conf.Server().JwtSignKey)
}

// newDB
func newDB(conf *config.Config) *mysql.GroupManager{
	return mysql.New(conf.Mysql())
}

// newRedis
func newRedis(conf *config.Config) *redis.Client{
	return redis.New(conf.Redis())
}

func newEtcd(conf *config.Config) *etcd.Client  {
	return etcd.New(conf.Etcd())
}