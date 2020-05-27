package http

import (
	"fmt"
	"github.com/ebar-go/ego/config"
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

	tree tree
}

type tree struct {
	nodes map[string]node
}

type node struct {
	name string
	path string
	children []node
}

// NewServer 实例化server
func NewServer() *Server {
	router := gin.Default()

	// use global trace middleware
	router.Use(middleware.Trace)

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
	// use lock
	server.mu.Lock()

	if len(args) > 1 {
		return fmt.Errorf("args must be less than one")
	}

	// 解析port
	port := getPortFromArgs(args...)

	// 404
	server.Router.NoRoute(server.NotFoundHandler)
	server.Router.NoMethod(server.NotFoundHandler)

	completeHost := net.JoinHostPort("", strconv.Itoa(port))

	return server.Router.Run(completeHost)
}

// getPortFromArgs
func getPortFromArgs(args ...int) int {
	port := config.Server().Port
	if len(args) == 1 {
		port = args[0]
		config.Server().Port = port
	}
	return port
}