package http

import (
	"errors"
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/ego/http/handler"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
	"sync"
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

// Start 启动服务
// port http server port
func (server *Server) Start(args ...int) error {
	port := app.Config().ServicePort
	if len(args) == 1 {
		port = args[0]
		app.Config().ServicePort = port
	}else if len(args) > 1 {

		return errors.New("length of args must less than 1")
	}

	// 防重复操作
	server.mu.Lock()

	// 404
	server.Router.NoRoute(server.NotFoundHandler)
	server.Router.NoMethod(server.NotFoundHandler)

	completeHost := net.JoinHostPort("", strconv.Itoa(port))

	return server.Router.Run(completeHost)
}
