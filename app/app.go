package app

import (
	"github.com/ebar-go/ego/component/auth"
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/etcd"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/mysql"
	"github.com/ebar-go/ego/component/redis"
	"github.com/ebar-go/egu"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/dig"
	"net/http"
)

var container *dig.Container

func init()  {
	container = dig.New()
	// 注入配置文件
	_ = container.Provide(config.New)
	// 注入http客户端
	_ = container.Provide(newHttpClient)
	// 注入日志管理器
	_ = container.Provide(newLogger)
	// 注入jwt组件
	_ = container.Provide(newJwt)
	// 注入bufferPool
	_ = container.Provide(egu.NewBufferPool)
	// 注入redis组件
	_ = container.Provide(newRedis)
	// 注入DB组件
	_ = container.Provide(newDB)
	// 注入etcd主键
	_ = container.Provide(newEtcd)


}

// Container 容器
func Container() *dig.Container {
	return container
}

// Config 配置文件
func Config() (conf *config.Config)  {
	_ = container.Invoke(func(c *config.Config) {
		conf = c
	})
	return
}

// Redis get redis connection
func Redis() (client *redis.Client) {
	_ = container.Invoke(func(cli *redis.Client) {
		client = cli
	})
	return
}

// InitDB 初始化DB
func InitDB() error  {
	 return  container.Invoke(func(gm *mysql.GroupManager) error {
		return gm.Connect()
	})
}

// DB 返回数据库连接
func DB() *gorm.DB {
	return GetDB("default")
}

// GetDB 通过名称获取数据库连接
func GetDB(name string) (conn *gorm.DB) {
	_ = container.Invoke(func(gm *mysql.GroupManager) {
		conn = gm.GetConnection(name)
	})
	return
}

// Http client
func Http() (client *http.Client) {
	_ = container.Invoke(func(instance *http.Client) {
		client = instance
	})
	return
}


// Logger 日志管理器
func Logger() (logger *log.Logger) {
	_ = container.Invoke(func(instance *log.Logger) {
		logger = instance
	})
	return
}

// Jwt jwt组件
func Jwt() (jwt *auth.JwtAuth) {
	_ = container.Invoke(func(instance *auth.JwtAuth) {
		jwt = instance
	})
	return
}

// Etcd
func Etcd() (client *etcd.Client)  {
	_ = container.Invoke(func(instance *etcd.Client) {
		client = instance
	})
	return
}

func BufferPool() (pool *egu.BufferPool) {
	_ = container.Invoke(func(instance *egu.BufferPool) {
		pool = instance
	})
	return
}