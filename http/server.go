package http

import (
	"github.com/ebar-go/ego/http/handler"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/event"
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
	"sync"
)

// Server Web服务管理器
type Server struct {
	// 并发锁
	mu sync.Mutex

	// gin的路由
	Router *gin.Engine

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
func (server *Server) Start(port int) error {
	// 防重复操作
	server.mu.Lock()

	// 准备容器
	event.PrepareContainer()

	// 准备日志管理器
	event.PrepareLogManager()

	// 404
	server.Router.NoRoute(server.NotFoundHandler)
	server.Router.NoMethod(server.NotFoundHandler)

	completeHost := net.JoinHostPort("", strconv.Itoa(port))

	return server.Router.Run(completeHost)
}
