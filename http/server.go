package http

import (
	"errors"

	"github.com/ebar-go/ego/http/constant"
	"github.com/ebar-go/ego/http/handler"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/log"
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
)

// Server Web服务管理器
type Server struct {
	LogPath         string // 日志路径
	AppDebug        bool
	SystemName      string      // 系统名称
	Address         string      // 启动地址,如果为空,默认是127.0.0.1
	Port            int         // 端口号
	Router          *gin.Engine // gin的路由
	initialize      bool
	JwtKey          []byte                 // jwt秘钥
	AllowCors       bool                   // 是否允许跨域
	NotFoundHandler func(ctx *gin.Context) // 404的处理方法
	Recover         func(ctx *gin.Context) // 接受panic的recover处理方法
}

// Init 服务初始化
func (server *Server) Init() error {
	if server.initialize {
		return errors.New("请勿重复初始化Http Server")
	}

	if server.Port == 0 {
		return errors.New("端口号不能为0")
	}

	server.Router = gin.Default()

	server.loadMiddleware()

	// 404
	server.Router.NoRoute(server.NotFoundHandler)
	server.Router.NoMethod(server.NotFoundHandler)

	server.initLogger()
	server.initialize = true
	return nil
}

// loadMiddleware 加载中间件
func (server *Server) loadMiddleware() {
	if server.Recover == nil {
		server.Recover = middleware.Recover
	}
	// recover
	server.Router.Use(server.Recover)
	// 请求日志
	server.Router.Use(middleware.RequestLog)

	if server.AllowCors {
		server.Router.Use(middleware.CORS)
	}

	middleware.SetJwtSigningKey(server.JwtKey)

	if server.NotFoundHandler == nil {
		server.NotFoundHandler = handler.NotFoundHandler
	}
}

// initLogger 初始化日志管理器
func (server *Server) initLogger() {
	// 初始化日志目录
	if server.LogPath == "" {
		server.LogPath = constant.DefaultLogPath
	}

	if server.SystemName == "" {
		server.SystemName = constant.DefaultSystemName
	}

	// 初始化系统日志管理器
	log.InitManager(log.ManagerConf{
		SystemName: server.SystemName,
		SystemPort: server.Port,
		LogPath:    server.LogPath,
		AppDebug:   server.AppDebug,
	})
}

// GetCompleteHost 获取完整的host
func (server Server) GetCompleteHost() string {
	return net.JoinHostPort(server.Address, strconv.Itoa(server.Port))
}

// Start 启动服务
func (server *Server) Start() error {
	if !server.initialize {
		return errors.New("请先初始化服务")
	}

	return server.Router.Run(server.GetCompleteHost())
}
