package http

import (
	"errors"
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/ego/event"
	"github.com/ebar-go/ego/http/handler"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
	"sync"
)


const (
	// mysql connect event
	MySqlConnectEvent = "MYSQL_CONNECT_EVENT"

	// redis connect event
	RedisConnectEvent = "REDIS_CONNECT_EVENT"
)

// Server Web服务管理器
type Server struct {
	// 并发锁,可导出结构体采用私有变量，而不采用内嵌的方式
	mu sync.Mutex

	// gin的路由
	Router *gin.Engine

	// not found handler
	NotFoundHandler func(ctx *gin.Context)
}

func init() {
	// register before start events
	app.EventDispatcher().AddListener(MySqlConnectEvent,
		event.NewListener(func(ev event.Event) {
			app.Mysql()
		}))

	app.EventDispatcher().AddListener(RedisConnectEvent,
		event.NewListener(func(ev event.Event) {
			app.Redis()
		}))
}

// NewServer 实例化server
func NewServer() *Server {
	router := gin.Default()

	// 默认引入请求日志中间件
	router.Use(middleware.RequestLog)

	return &Server{
		Router:          router,
		NotFoundHandler: handler.NotFoundHandler,
	}
}

// Start run http server
// args must be less than one,
// eg: Start()
//	   Start(8080)
func (server *Server) Start(args ...int) error {
	port := app.Config().ServicePort
	if len(args) == 1 {
		port = args[0]
		app.Config().ServicePort = port
	} else if len(args) > 1 {
		return errors.New("args must be less than one")
	}

	// 防重复操作
	server.mu.Lock()

	// 404
	server.Router.NoRoute(server.NotFoundHandler)
	server.Router.NoMethod(server.NotFoundHandler)

	completeHost := net.JoinHostPort("", strconv.Itoa(port))

	// mysql auto connect
	if app.Config().Mysql().AutoConnect {
		app.EventDispatcher().Trigger(MySqlConnectEvent, nil)
	}

	// redis auto connect
	if app.Config().Redis().AutoConnect {
		app.EventDispatcher().Trigger(RedisConnectEvent, nil)
	}

	return server.Router.Run(completeHost)
}

