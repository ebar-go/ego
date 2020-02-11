package app

import (
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/utils"
	"github.com/ebar-go/event"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"time"
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
)

func init() {
	// init event dispatcher
	utils.FatalError("InitEventDispatcher", Container.Provide(event.NewDispatcher))

	// use eventDispatcher manage global service initialize
	eventDispatcher := EventDispatcher()

	eventDispatcher.Register(LogManagerInitEvent, func(ev event.Event) {
		utils.FatalError("InitLogManager", initLogManager())
	})

	eventDispatcher.Register(MySqlConnectEvent, func(ev event.Event) {
		utils.FatalError("ConnectDatabase", connectDatabase())
	})

	eventDispatcher.Register(RedisConnectEvent, func(ev event.Event) {
		utils.FatalError("ConnectRedis", connectRedis())
	})

}

// initLogManager
func initLogManager() error {
	return Container.Provide(func() log.Manager {
		return log.NewManager(log.ManagerConf{
			SystemName: config.Server().Name,
			SystemPort: config.Server().Port,
			LogPath:    config.Server().LogPath,
		})
	})
}

// connectRedis
func connectRedis() error {
	return Container.Provide(func() (*redis.Client, error) {
		connection := redis.NewClient(config.Redis().Options())
		_, err := connection.Ping().Result()
		if err != nil {
			return nil, errors.RedisConnectFailed("%s", err.Error())
		}

		return connection, nil
	})
}

// connectDatabase
func connectDatabase() error {
	return Container.Provide(func() (*gorm.DB, error) {
		options := config.Mysql()
		connection, err := gorm.Open("mysql", options.Dsn())
		if err != nil {
			return nil, errors.MysqlConnectFailed("%s", err.Error())
		}

		// set log mod
		connection.LogMode(options.LogMode)
		// set pool config
		connection.DB().SetMaxIdleConns(options.MaxIdleConnections)
		connection.DB().SetMaxOpenConns(options.MaxOpenConnections)
		connection.DB().SetConnMaxLifetime(time.Duration(options.MaxLifeTime) * time.Second)

		return connection, nil
	})
}

