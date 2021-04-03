package app

import (
	"github.com/ebar-go/ego/component/auth"
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/etcd"
	selflog "github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/mysql"
	"github.com/ebar-go/ego/component/redis"
	"github.com/ebar-go/ego/http"
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
	return  &App{container: buildContainer()}
}
// Container 容器
func (app *App) Container() *dig.Container {
	return app.container
}
// buildContainer 构造容器
func buildContainer() *dig.Container {
	container := dig.New()
	// 注入配置项组件
	config.Inject(container)

	// 日志
	_ = container.Provide(selflog.New)

	// 连接数据库
	_ = container.Provide(mysql.Connect)

	// redis服务
	_ = container.Provide(redis.Connect)

	// http server
	http.Inject(container)

	// jwt
	_ = container.Provide(func(config *http.Config) auth.Jwt{
		return auth.NewJwt(config.JwtSignKey)
	})

	// etcd
	_ = container.Provide(etcd.New)

	// 定时任务
	_ = container.Provide(cron.New)
	return container
}

// LoadConfig 加载配置文件
func (app *App) LoadConfig(path ...string) error {
	return app.container.Invoke(func(config *config.Config) error{
		return config.LoadFile(path...)
	})
}
// LoadRouter 加载路由
func (app *App) LoadRouter(loader interface{}) error {
	return app.container.Invoke(loader)
}
// LoadTask 加载定时任务
func (app *App) LoadTask(loader func(cron *cron.Cron)) error {
	return app.container.Invoke(loader)
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

// StartCron 开启定时任务
func (app *App) StartCron() error {
	return app.container.Invoke(func(cron *cron.Cron){
		cron.Start()
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

	// close task
	_ = app.container.Invoke(func(c *cron.Cron) {
		c.Stop()
	})
}

