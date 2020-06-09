package app

import (
	"fmt"
	"github.com/ebar-go/ego/component/event"
	"github.com/ebar-go/ego/component/mysql"
	"github.com/ebar-go/ego/config"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"time"
)


// InitRedis 初始化redis
func InitRedis() error {
	connection := redis.NewClient(config.Redis().Options())
	_, err := connection.Ping().Result()
	if err != nil {
		return err
	}

	return Container.Provide(func() *redis.Client {
		return connection
	})
}
// InitDB 初始化DB
func InitDB() error {
	dialect := "mysql"
	group := config.MysqlGroup()
	if group.Items == nil {
		return fmt.Errorf("mysql config is empty")
	}

	for name, item := range group.Items {
		dataSourceItems := item.DsnItems()

		adapter, err := mysql.NewReadWriteAdapter(dialect, dataSourceItems)
		if err != nil {
			return err
		}

		adapter.SetMaxIdleConns(item.MaxIdleConnections)
		adapter.SetMaxOpenConns(item.MaxOpenConnections)
		adapter.SetConnMaxLifetime(time.Duration(item.MaxLifeTime) * time.Second)

		conn, err := gorm.Open(dialect, adapter)
		if err != nil {
			return err
		}
		dbGroup[name] = conn
	}

	event.Trigger(event.AfterDatabaseConnect, nil)

	return nil
}
