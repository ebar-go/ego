package http

import (
	"errors"

	"github.com/ebar-go/ego/http/handler"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/gin-gonic/gin"
)

// Server Web服务管理器
type Server struct {
	Address    string
	Router     *gin.Engine
	initialize bool
}

// Init 服务初始化
func (server *Server) Init() error {
	if server.initialize {
		return errors.New("请勿重复初始化Http Server")
	}

	server.Router = gin.Default()

	// 请求日志
	server.Router.Use(middleware.RequestLog)

	// 404
	server.Router.NoRoute(handler.NotFoundHandler)
	server.Router.NoMethod(handler.NotFoundHandler)

	server.initialize = true
	return nil
}

// Start 启动服务
func (server *Server) Start() error {
	if !server.initialize {
		return errors.New("请先初始化服务")
	}

	return server.Router.Run(server.Address)
}
