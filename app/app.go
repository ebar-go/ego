package app

import (
	"github.com/ebar-go/ego/component/auth"
	"github.com/ebar-go/ego/component/buffer"
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/dig"
	"net/http"
)

var (
	Container = dig.New()
	dbGroup   = make(map[string]*gorm.DB)
)

func init()  {
	// 注入配置文件
	_ = Container.Provide(config.New)
	// 注入http客户端
	_ = Container.Provide(newHttpClient)
	// 注入日志管理器
	_ = Container.Provide(newLogger)
	// 注入jwt组件
	_ = Container.Provide(func(conf *config.Config) *auth.JwtAuth{
		return auth.New(conf.Server().JwtSignKey)
	})
	_ = Container.Provide(func(conf *config.Config) *redis.Redis {
		return &redis.Redis{
			Options: conf.Redis().Options(),
		}
	})
	_ = Container.Provide(buffer.NewPool)
}

// Config 配置文件
func Config() (conf *config.Config)  {
	_ = Container.Invoke(func(c *config.Config) {
		conf = c
	})
	return
}

// Redis get redis connection
func Redis() (connection *redis.Redis) {
	_ = Container.Invoke(func(conn *redis.Redis) {
		connection = conn
	})
	return
}

// DB 返回数据库连接
func DB() *gorm.DB {
	return dbGroup["default"]
}

// GetDB 通过名称获取数据库连接
func GetDB(connectionName string) *gorm.DB {
	return dbGroup[connectionName]
}

// Http client
func Http() (client *http.Client) {
	_ = Container.Invoke(func(instance *http.Client) {
		client = instance
	})
	return
}


// Logger 日志管理器
func Logger() (logger *log.Logger) {
	_ = Container.Invoke(func(instance *log.Logger) {
		logger = instance
	})
	return
}

// Jwt jwt组件
func Jwt() (jwt *auth.JwtAuth) {
	_ = Container.Invoke(func(instance *auth.JwtAuth) {
		jwt = instance
	})
	return
}


// BufferPool buffer池
func BufferPool() (pool *buffer.Pool) {
	_ = Container.Invoke(func(instance *buffer.Pool) {
		pool = instance
	})
	return
}
