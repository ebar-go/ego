package app

import (
	"fmt"
	"github.com/ebar-go/ego/component/auth"
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/etcd"
	"github.com/ebar-go/ego/component/mysql"
	"github.com/ebar-go/ego/component/redis"
	"github.com/ebar-go/ego/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"log"
)

type App struct {
	container *dig.Container
}

func New() *App {
	app := &App{container: dig.New()}
	if err := app.inject(); err != nil {
		log.Fatalf("%v\n", err)
	}
	return app
}

func (app *App) Container() *dig.Container {
	return app.container
}

func (app *App) inject() error {
	if err := app.container.Provide(config.New); err != nil {
		log.Printf("inject config: %v\n", err)
	}

	// 数据库配置
	if err := app.container.Provide(newDatabaseConfig); err != nil {
		log.Printf("inject database config: %v\n", err)
	}

	// redis配置
	if err := app.container.Provide(newRedisConfig); err != nil {
		log.Printf("inject database config: %v\n", err)
	}

	// etcd config
	if err := app.container.Provide(newEtcdConfig); err != nil {
		log.Printf("inject etcd config: %v\n", err)
	}

	if err := app.container.Provide(newLogger); err != nil {
		log.Printf("inject logger: %v\n", err)
	}

	// 连接数据库
	if err := app.container.Provide(mysql.Connect); err != nil {
		log.Printf("inject database: %v\n", err)
	}

	// redis服务
	if err := app.container.Provide(redis.Connect); err != nil {
		return fmt.Errorf("connect redis: %v", err)
	}

	// http server
	if err := app.container.Provide(http.New); err != nil {
		return err
	}

	// router
	if err := app.container.Provide(func(server *http.Server) *gin.Engine{
		return server.Router()
	}); err != nil {
		return err
	}

	// jwt
	if err := app.container.Provide(func(config *config.Config) auth.Jwt{
		return auth.NewJwt(config.JwtSignKey)
	}); err != nil {
		log.Printf("inject jwt: %v\n", err)
	}

	// etcd
	if err  := app.container.Provide(etcd.New); err != nil {
		log.Printf("inject etcd: %v\n", err)
	}
	return nil
}

func (app *App) ServeHttp() error {
	return app.container.Invoke(func(server *http.Server) error {
		return server.Serve()
	})
}

func (app *App) LoadConfig(path ...string) error {
	return app.container.Invoke(func(config *config.Config) error{
		return config.LoadFile(path...)
	})
}

func (app *App) ServeRpc() error {
	return nil
}

func (app *App) ServerWebsocket() error {
	return nil
}


