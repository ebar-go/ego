package http

import (
	"github.com/gin-gonic/gin"
	"errors"
	"github.com/ebar-go/ego/log"
	"github.com/ebar-go/ego/http/middleware"
)

type Server struct {
	Address string
	Router *gin.Engine
	initialize bool
	SystemLogHandler *log.Logger
}

// Init 服务初始化
func (server *Server)Init() error {
	if server.initialize {
		return errors.New("请勿重复初始化Http Server")
	}
	server.Router = gin.New()

	server.Router.Use(middleware.RequestLog)

	if server.SystemLogHandler == nil {
		server.SystemLogHandler = log.DefaultLogger()
	}

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
