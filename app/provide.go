package app

import (
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/mns"
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/event"
	"github.com/ebar-go/ego/utils"
	"github.com/ebar-go/ego/ws"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron"
)

const (
	// config init event
	ConfigInitEvent = "CONFIG_INIT_EVENT"

	// log manager init event
	LogManagerInitEvent = "LOG_MANAGER_INIT_EVENT"

	// mysql connect event
	MySqlConnectEvent = "MYSQL_CONNECT_EVENT"

	// redis connect event
	RedisConnectEvent = "REDIS_CONNECT_EVENT"

	// websocket init event
	WebSocketInitEvent = "WS_INIT_EVENT"

	// mns client init event
	MNSClientInitEvent = "MNS_INIT_EVENT"

	// task manager init event
	TaskManagerInitEvent = "TASK_INIT_EVENT"

)

func init() {
	// init event dispatcher
	utils.FatalError("InitEventDispatcher", initEventDispatcher())

	// use eventDispatcher manage global service initialize
	eventDispatcher := EventDispatcher()

	eventDispatcher.Register(ConfigInitEvent, func(ev event.Event) {
		utils.FatalError("InitConfig", initConfig())
	})

	eventDispatcher.Register(LogManagerInitEvent, func(ev event.Event) {
		utils.FatalError("InitLogManager", initLogManager())
	})

	eventDispatcher.Register(MySqlConnectEvent, func(ev event.Event) {
		utils.FatalError("ConnectDatabase", connectDatabase())
	})


	eventDispatcher.Register(MySqlConnectEvent, func(ev event.Event) {
		utils.FatalError("ConnectDatabase", connectDatabase())
	})

	eventDispatcher.Register(RedisConnectEvent, func(ev event.Event) {
		utils.FatalError("ConnectRedis", connectRedis())
	})

	eventDispatcher.Register(WebSocketInitEvent, func(ev event.Event) {
		utils.LogError("InitWebSocketManager", initWebSocketManager())
	})

	eventDispatcher.Register(MNSClientInitEvent, func(ev event.Event) {
		utils.LogError("InitMnsClient", initMnsClient())
	})

	eventDispatcher.Register(TaskManagerInitEvent, func(ev event.Event) {
		utils.LogError("InitTaskManager", initTaskManager())
	})

}

// initConfig
func initConfig() error {
	return app.Provide(config.NewInstance)
}

// initLogManager
func initLogManager() error {
	return app.Provide(func(conf *config.Config) log.Manager {
		return log.NewManager(log.ManagerConf{
			SystemName: conf.ServiceName,
			SystemPort: conf.ServicePort,
			LogPath:    conf.LogPath,
		})
	})
}

// initTaskManager
func initTaskManager() error {
	return app.Provide(cron.New)
}

// initWebSocketManager
func initWebSocketManager() error {
	return app.Provide(ws.NewManager)
}

// connectRedis
func connectRedis() error {
	return app.Provide(func(conf *config.Config) (*redis.Client, error) {
		connection := redis.NewClient(conf.Redis().Options())
		_, err := connection.Ping().Result()
		if err != nil {
			return nil, errors.RedisConnectFailed("%s", err.Error())
		}

		return connection, nil
	})
}

// connectDatabase
func connectDatabase() error {
	return app.Provide(func(conf *config.Config) (*gorm.DB, error) {
		options := conf.Mysql()
		connection, err := gorm.Open("mysql", options.Dsn())
		if err != nil {
			return nil, errors.MysqlConnectFailed("%s", err.Error())
		}

		// set log mod
		connection.LogMode(options.LogMode)
		// set pool config
		connection.DB().SetMaxIdleConns(options.MaxIdleConnections)
		connection.DB().SetMaxOpenConns(options.MaxOpenConnections)

		return connection, nil
	})
}

// initMnsClient
func initMnsClient() error {
	return app.Provide(func(conf *config.Config, logManager log.Manager) (mns.Client) {
		mnsConfig := conf.Mns()
		return mns.NewClient(
			mnsConfig.Url,
			mnsConfig.AccessKeyId,
			mnsConfig.AccessKeySecret,
			logManager)
	})
}

// initEventDispatcher
func initEventDispatcher() error {
	return app.Provide(event.NewDispatcher)
}
