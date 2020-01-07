package app

import (
	"github.com/ebar-go/ego/component/auth"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/mns"
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/event"
	"github.com/ebar-go/ego/utils"
	"github.com/ebar-go/ego/ws"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/robfig/cron"
	"go.uber.org/dig"
	"sync"
)

var (
	app = NewContainer()
	httpClientInit sync.Once
)

// NewContainer return an empty container
func NewContainer() *dig.Container {
	return dig.New()
}


// Config return Config instance
func Config() (conf *config.Config) {
	if err := app.Invoke(func(c *config.Config) {
		conf = c
	}); err != nil {
		// use sync.once is better ?
		conf = config.NewInstance()

		// if get err , panic
		utils.PanicErr("InitConfig", app.Provide(func() *config.Config {
			return conf
		}))
	}

	return
}

// LogManager return log manager
func LogManager() (manager log.Manager) {
	if err := app.Invoke(func(m log.Manager) {
		manager = m
	}); err != nil {
		conf := Config()
		manager = log.NewManager(log.ManagerConf{
			SystemName: conf.ServiceName,
			SystemPort: conf.ServicePort,
			LogPath:    conf.LogPath,
		})

		utils.PanicErr("InitLogManager", app.Provide(func() log.Manager {
			return manager
		}))
	}

	return
}

// Task return task manager
func Task() (manager *cron.Cron) {
	if err := app.Invoke(func(c *cron.Cron) {
		manager = c
	}); err != nil {
		manager = cron.New()

		utils.PanicErr("InitTaskManager", app.Provide(func() *cron.Cron {
			return manager
		}))
	}

	return
}

// Jwt return a jwt instance
func Jwt() (jwt auth.Jwt) {
	if err := app.Invoke(func(j auth.Jwt) {
		jwt = j
	}); err != nil {
		jwt = auth.New(Config().JwtSignKey)
		_ = app.Provide(func() auth.Jwt {
			return jwt
		})
	}

	return
}

// WebSocket return ws manager
func WebSocket() (manager ws.Manager) {
	if err := app.Invoke(func(m ws.Manager) {
		manager = m
	}); err != nil {
		manager = ws.NewManager()
		utils.PanicErr("InitWebSocketManager", app.Provide(func() ws.Manager {
			return manager
		}))
	}

	return
}

// Redis get redis connection
func Redis() (connection *redis.Client) {
	if err := app.Invoke(func(conn *redis.Client) {
		connection = conn
	}); err != nil {
		connection = redis.NewClient(Config().Redis().Options())
		_, err = connection.Ping().Result()
		if err != nil {
			panic(errors.RedisConnectFailed("%s", err.Error()))
		}
		_ = app.Provide(func() *redis.Client {
			return connection
		})
	}

	return connection
}

// Mysql return mysql connection
func Mysql() (connection *gorm.DB) {
	if err := app.Invoke(func(conn *gorm.DB) {
		connection = conn
	}); err != nil {
		conf := Config().Mysql()

		connection, err = gorm.Open("mysql", conf.Dsn())
		if err != nil {
			panic(errors.MysqlConnectFailed("%s", err.Error()))
		}

		// set log mod
		connection.LogMode(conf.LogMode)
		// set pool config
		connection.DB().SetMaxIdleConns(conf.MaxIdleConnections)
		connection.DB().SetMaxOpenConns(conf.MaxOpenConnections)

		_ = app.Provide(func() *gorm.DB {
			return connection
		})
	}

	return connection
}

// Mns return ali yun mns client
func Mns() (client mns.Client) {
	if err := app.Invoke(func(cli mns.Client) {
		client = cli
	}); err != nil {
		conf := Config().Mns()
		client = mns.NewClient(conf.Url, conf.AccessKeyId, conf.AccessKeySecret, LogManager())
		utils.PanicErr("InitMns", app.Provide(func() mns.Client {
			return client
		}))
	}

	return client
}

// EventDispatcher get event dispatcher instance
func EventDispatcher() (dispatcher event.Dispatcher) {
	if err := app.Invoke(func(d event.Dispatcher) {
		dispatcher = d
	}); err != nil {
		dispatcher = event.NewDispatcher()

		utils.PanicErr("InitEventDispatcher", app.Provide(func() event.Dispatcher {
			return dispatcher
		}))
	}
	return
}
