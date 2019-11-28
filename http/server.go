package http

import (
	"github.com/ebar-go/ego/http/handler"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/log"
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
	"sync"
)


const (
	defaultPort = 8080
	defaultLogPath = "/tmp/log"
	defaultName = "app"
)

// Server Web服务管理器
type Server struct {
	initialize *sync.Mutex
	// 系统名称，可选
	name string

	// 启动地址,如果为空,默认是0.0.0.0
	address string

	// 端口号,默认是8080
	port int

	// 日志路径
	logPath string

	// 是否设置app日志等级为debug
	appDebug bool


	Router *gin.Engine // gin的路由


	jwtKey []byte // jwt秘钥


	// 404的处理方法
	notFoundHandler func(ctx *gin.Context)

	// 接受panic的recover处理方法
	recoverHandler func(ctx *gin.Context)
}

// 实例化server
func NewServer() *Server {
	router := gin.Default()

	// 默认引入请求日志中间件
	router.Use(middleware.RequestLog)

	return &Server{
		name: defaultName,
		port: defaultPort,
		notFoundHandler: handler.NotFoundHandler,
		appDebug: false,
		Router: router,
		initialize:new(sync.Mutex),
		logPath: defaultLogPath,
	}
}

// SetName 设置系统名称
func (server *Server) SetName(name string) {
	server.name = name
}

// SetLogPath 设置日志路径
func (server *Server) SetLogPath(path string) {
	server.logPath = path
}

// AppDebug 是否开启app日志的debug等级
func (server *Server) AppDebug(debug bool) {
	server.appDebug = debug
}

// SetJwtKey 设置jwt
func (server *Server) SetJwtKey(key []byte) {
	server.jwtKey = key
}

// SetNotFoundHandler 设置404处理器
func (server *Server) SetNotFoundHandler(notFoundHandler func(ctx *gin.Context)) {
	server.notFoundHandler = notFoundHandler
}

// SetAddress 设置地址
func (server *Server) SetAddress(address string) {
	server.address = address
}

// SetPort 设置端口
func (server *Server) SetPort(port int) {
	server.port = port
}

// Start 启动服务
func (server *Server) Start() error {
	server.initialize.Lock()

	// 404
	server.Router.NoRoute(server.notFoundHandler)
	server.Router.NoMethod(server.notFoundHandler)

	middleware.SetJwtSigningKey(server.jwtKey)

	// 初始化系统日志管理器
	log.InitManager(log.ManagerConf{
		SystemName: server.name,
		SystemPort: server.port,
		LogPath: server.logPath,
		AppDebug: server.appDebug,
	})

	completeHost := net.JoinHostPort(server.address, strconv.Itoa(server.port))

	return server.Router.Run(completeHost)
}
