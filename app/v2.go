package app

import (
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/mysql"
	"github.com/ebar-go/ego/component/redis"
	"github.com/ebar-go/ego/http/server"
	"go.uber.org/dig"
	"log"
)

type App struct {
	container *dig.Container
}

func (app *App) Container() *dig.Container {
	return app.container
}

func (app *App) Run() error {
	if err := app.container.Provide(config.New); err != nil {
		log.Printf("inject config: %v\n", err)
	}

	if err := app.container.Provide(server.NewConfig); err != nil {
		log.Printf("inject server config: %v\n", err)
	}

	// 数据库配置
	if err := app.container.Provide(func (config *config.Config) *mysql.Config{
		return &mysql.Config{}
	}); err != nil {
		log.Printf("inject database config: %v\n", err)
	}

	// 连接数据库
	if err := app.container.Provide(mysql.Connect); err != nil {
		log.Printf("inject database: %v\n", err)
	}

	// redis配置

	// redis服务
	if err := app.container.Provide(redis.Connect); err != nil {
		return err
	}

	return nil
}

