package app

import (
	"github.com/ebar-go/ego/component/auth"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/mns"
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ego/event"
	"github.com/ebar-go/ego/utils"
	"github.com/ebar-go/ego/ws"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/robfig/cron"
	"go.uber.org/dig"
)

var (
	app = NewContainer()
)

// NewContainer return an empty container
func NewContainer() *dig.Container {
	return dig.New()
}

// Config return Config instance
func Config() (conf *config.Config) {
	_ = app.Invoke(func(c *config.Config) {
		conf = c
	})
	return
}

// LogManager return log manager
func LogManager() (manager log.Manager) {
	_ = app.Invoke(func(m log.Manager) {
		manager = m
	})

	return
}

// Task return task manager
func Task() (manager *cron.Cron) {
	_ = app.Invoke(func(c *cron.Cron) {
		manager = c
	})

	return
}

// Jwt return a jwt instance
func Jwt() (jwt auth.Jwt) {
	_ = app.Invoke(func(j auth.Jwt) {
		jwt = j
	})

	return
}

// WebSocket return ws manager
func WebSocket() (manager ws.Manager) {
	if err := app.Invoke(func(m ws.Manager) {
		manager = m
	}); err != nil {
		manager = ws.NewManager()
		utils.FatalError("InitWebSocketManager", app.Provide(func() ws.Manager {
			return manager
		}))
	}

	return
}

// Redis get redis connection
func Redis() (connection *redis.Client) {
	_ = app.Invoke(func(conn *redis.Client) {
		connection = conn
	})
	return
}

// Mysql return mysql connection
func Mysql() (connection *gorm.DB) {
	_ = app.Invoke(func(conn *gorm.DB) {
		connection = conn
	})
	return
}

// Mns return ali yun mns client
func Mns() (client mns.Client) {
	_ =  app.Invoke(func(cli mns.Client) {
		client = cli
	})
	return
}

// EventDispatcher get event dispatcher instance
func EventDispatcher() (dispatcher event.Dispatcher) {
	_ = app.Invoke(func(d event.Dispatcher) {
		dispatcher = d
	})
	return
}
