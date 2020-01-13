package app

import (
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/mns"
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ego/event"
	"github.com/ebar-go/ego/ws"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/robfig/cron"
	"go.uber.org/dig"
)

var (
	DefaultContainer = NewContainer()
)

// NewContainer return an empty container
func NewContainer() *dig.Container {
	return dig.New()
}

// Config return Config instance
func Config() (conf config.Config) {
	_ = DefaultContainer.Invoke(func(c config.Config) {
		conf = c
	})
	return
}

// LogManager return log manager
func LogManager() (manager log.Manager) {
	_ = DefaultContainer.Invoke(func(m log.Manager) {
		manager = m
	})

	return
}

// Task return task manager
func Task() (manager *cron.Cron) {
	_ = DefaultContainer.Invoke(func(c *cron.Cron) {
		manager = c
	})

	return
}

// WebSocket return ws manager
func WebSocket() (manager ws.Manager) {
	_ = DefaultContainer.Invoke(func(m ws.Manager) {
		manager = m
	})
	return
}

// Redis get redis connection
func Redis() (connection *redis.Client) {
	_ = DefaultContainer.Invoke(func(conn *redis.Client) {
		connection = conn
	})
	return
}

// Mysql return mysql connection
func Mysql() (connection *gorm.DB) {
	_ = DefaultContainer.Invoke(func(conn *gorm.DB) {
		connection = conn
	})
	return
}

// Mns return ali yun mns client
func Mns() (client mns.Client) {
	_ =  DefaultContainer.Invoke(func(cli mns.Client) {
		client = cli
	})
	return
}

// EventDispatcher get event dispatcher instance
func EventDispatcher() (dispatcher event.Dispatcher) {
	_ = DefaultContainer.Invoke(func(d event.Dispatcher) {
		dispatcher = d
	})
	return
}
