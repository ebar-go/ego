package app

import (
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/mns"
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ws"
	"github.com/ebar-go/event"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/dig"
)

var (
	Container = NewContainer()
)

// NewContainer return an empty container
func NewContainer() *dig.Container {
	return dig.New()
}

// LogManager return log manager
func LogManager() (manager log.Manager) {
	_ = Container.Invoke(func(m log.Manager) {
		manager = m
	})

	return
}

// WebSocket return ws manager
func WebSocket() (manager ws.Manager) {
	if err := Container.Invoke(func(m ws.Manager) {
		manager = m
	}); err != nil {
		manager = ws.New()
		_ = Container.Provide(func() ws.Manager{
			return manager
		})
	}
	return
}

// Redis get redis connection
func Redis() (connection *redis.Client) {
	_ = Container.Invoke(func(conn *redis.Client) {
		connection = conn
	})
	return
}

// Mysql return mysql connection
func Mysql() (connection *gorm.DB) {
	_ = Container.Invoke(func(conn *gorm.DB) {
		connection = conn
	})
	return
}

// Mns return ali yun mns client
func Mns() (client mns.Client) {
	if err :=  Container.Invoke(func(cli mns.Client) {
		client = cli
	}); err != nil {
		mnsConfig := config.Mns()
		client = mns.NewClient(
			mnsConfig.Url,
			mnsConfig.AccessKeyId,
			mnsConfig.AccessKeySecret,
			LogManager())
		_ = Container.Provide(func() (mns.Client) {
			return client
		})
	}
	return
}

// EventDispatcher get event dispatcher instance
func EventDispatcher() (dispatcher event.Dispatcher) {
	_ = Container.Invoke(func(d event.Dispatcher) {
		dispatcher = d
	})
	return
}
