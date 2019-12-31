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

// NewContainer 新容器
func NewContainer() *dig.Container {
	return dig.New()
}

// Config 配置项
func Config() (conf *config.Config) {
	if err := app.Invoke(func(c *config.Config) {
		conf = c
	}); err != nil {
		conf = config.NewInstance()
		utils.LogError("InitConfig", app.Provide(func() *config.Config {
			return conf
		}))
	}

	return
}

// LogManager 日志管理器
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

		utils.LogError("InitLogManager", app.Provide(func() log.Manager {
			return manager
		}))
	}

	return
}

// Task 任务管理器
func Task() (manager *cron.Cron) {
	if err := app.Invoke(func(c *cron.Cron) {
		manager = c
	}); err != nil {
		manager = cron.New()

		utils.LogError("InitTaskManager", app.Provide(func() *cron.Cron {
			return manager
		}))
	}

	return
}

// Jwt JsonWebToken
func Jwt() (jwt auth.Jwt) {
	if err := app.Invoke(func(j auth.Jwt) {
		jwt = j
	}); err != nil {
		jwt = auth.NewJwt(Config().JwtSignKey)
		utils.LogError("InitJwt", app.Provide(func() auth.Jwt {
			return jwt
		}))
	}

	return
}

// WebSocket
func WebSocket() (manager ws.Manager) {
	if err := app.Invoke(func(m ws.Manager) {
		manager = m
	}); err != nil {
		manager = ws.NewManager()
		utils.LogError("InitWebSocketManager", app.Provide(func() ws.Manager {
			return manager
		}))
	}

	return
}

// Redis 获取redis连接
func Redis() (connection *redis.Client) {
	if err := app.Invoke(func(conn *redis.Client) {
		connection = conn
	}); err != nil {
		connection = redis.NewClient(Config().Redis().Options())
		_, err = connection.Ping().Result()
		utils.FatalError("InitRedis", err)
		_ = app.Provide(func() *redis.Client {
			return connection
		})
	}

	return connection
}

// Mysql
func Mysql() (connection *gorm.DB) {
	if err := app.Invoke(func(conn *gorm.DB) {
		connection = conn
	}); err != nil {
		conf := Config().Mysql()

		connection, err = gorm.Open("mysql", conf.Dsn())
		utils.FatalError("InitMysql", err)

		// 设置是否打印日志
		connection.LogMode(conf.LogMode)
		// 设置连接池
		connection.DB().SetMaxIdleConns(conf.MaxIdleConnections)
		connection.DB().SetMaxOpenConns(conf.MaxOpenConnections)

		_ = app.Provide(func() *gorm.DB {
			return connection
		})
	}

	return connection
}

// Mns 阿里云mns
func Mns() (client mns.Client) {
	if err := app.Invoke(func(cli mns.Client) {
		client = cli
	}); err != nil {
		conf := Config().Mns()
		client = mns.NewClient(conf.Url, conf.AccessKeyId, conf.AccessKeySecret, LogManager())
		_ = app.Provide(func() mns.Client {
			return client
		})
	}

	return client
}

// EventDispatcher get event dispatcher instance
func EventDispatcher() (dispatcher event.Dispatcher) {
	if err := app.Invoke(func(d event.Dispatcher) {
		dispatcher = d
	}); err != nil {
		dispatcher = event.NewDispatcher()
		utils.LogError("InitEventDispatcher", app.Provide(func() event.Dispatcher {
			return dispatcher
		}))
	}
	return
}