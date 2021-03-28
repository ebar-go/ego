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
	"github.com/robfig/cron"
	"go.uber.org/dig"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// App 应用
type App struct {
	container *dig.Container
}

// New 实例化
func New() *App {
	app := &App{container: dig.New()}
	if err := app.injectComponents(); err != nil {
		log.Fatalf("%v\n", err)
	}
	return app
}
// Container 容器
func (app *App) Container() *dig.Container {
	return app.container
}
// injectComponents 注入组件
func (app *App) injectComponents() error {
	// 注入配置项组件
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

	// 日志
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
	if err := app.container.Provide(http.NewServer); err != nil {
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

	// 定时任务
	_ = app.container.Provide(cron.New)
	return nil
}

// LoadConfig 加载配置文件
func (app *App) LoadConfig(path ...string) error {
	return app.container.Invoke(func(config *config.Config) error{
		return config.LoadFile(path...)
	})
}

// ServeHTTP 启动http服务
func (app *App) ServeHTTP() {
	_ =  app.container.Invoke(func(server *http.Server) {
		server.Serve()
	})
}


func (app *App) serveRPC() error {
	return nil
}

func (app *App) serveWS() error {
	return nil
}

// 加载task
func (app *App) LoadTask(loader func (c *cron.Cron)) {
	_ = app.container.Invoke(func(c *cron.Cron, conf *config.Config) {
		loader(c)

		if conf.Task {
			c.Start()
		}
	})
}

// 启动应用
func (app *App) Run() {
	// wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	// close http
	_ = app.container.Invoke(func(server *http.Server) {
		server.Close()
	})

	// clise task
	_ = app.container.Invoke(func(c *cron.Cron) {
		c.Stop()
	})
}

