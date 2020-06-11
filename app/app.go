package app

import (
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ego/constant"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/dig"
	"net"
	"net/http"
	"time"
)

var (
	Container = NewContainer()
	dbGroup   = make(map[string]*gorm.DB)
)

// NewContainer return an empty container
func NewContainer() *dig.Container {
	return dig.New()
}

// Redis get redis connection
func Redis() (connection *redis.Client) {
	_ = Container.Invoke(func(conn *redis.Client) {
		connection = conn
	})
	return
}

// DB 返回数据库连接
func DB() *gorm.DB {
	return dbGroup[constant.MysqlDefaultConnection]
}

// GetDB 通过名称获取数据库连接
func GetDB(connectionName string) *gorm.DB {
	return dbGroup[connectionName]
}

// Http client
func Http() (client *http.Client) {
	err := Container.Invoke(func(cli *http.Client) {
		client = cli
	})
	if err != nil {
		client = &http.Client{
			Transport: &http.Transport{ // 配置连接池
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				IdleConnTimeout: time.Duration(config.Server().HttpRequestTimeOut) * time.Second,
			},
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       time.Duration(config.Server().HttpRequestTimeOut) * time.Second,
		}

		_ = Container.Provide(func() *http.Client {
			return client
		})
	}
	return
}
