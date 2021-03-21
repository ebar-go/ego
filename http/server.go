package http

import (
	"context"
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/ego/component/event"
	"github.com/ebar-go/ego/http/handler"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/validator"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// Server Web服务管理器
type Server struct {
	// 并发锁,可导出结构体采用私有变量，而不采用内嵌的方式
	mu sync.Mutex

	router *gin.Engine

	// gin的路由
	options *options
}

// HttpServer 获取Server示例
func New(opts ...Option) *Server {
	defaultOption := &options{
		port:           app.Config().Server().Port,
		notFoundHandler: handler.NotFoundHandler,
	}

	for _, opt := range opts {
		opt.apply(defaultOption)
	}
	router := gin.Default()

	// use global trace middleware
	router.Use(middleware.Trace)

	return &Server{
		options: defaultOption,
		router: router,
	}
}

func (server *Server) beforeStart()  {
	binding.Validator = new(validator.Validator)
	// before start
	event.Trigger(event.BeforeHttpStart, nil)
	// 404
	server.router.NoRoute(server.options.notFoundHandler)
	server.router.NoMethod(server.options.notFoundHandler)

	if app.Config().Server().Pprof {
		pprof.Register(server.router)
	}

	if app.Config().Server().Task {
		go app.Task().Start()
	}
}

// RouteLoader 加载路由
func (server *Server) RouteLoader(loader func (router *gin.Engine)) {
	loader(server.router)
}

// Start run http server
func (server *Server) Start() error {
	// use lock
	server.mu.Lock()

	server.beforeStart()

	completeHost := net.JoinHostPort("", strconv.Itoa(server.options.port))

	srv := &http.Server{
		Addr:    completeHost,
		Handler: server.router,
	}

	go func() {
		log.Printf("Listening and serving HTTP on %s\n", completeHost)

		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%s\n", err)
		}

		// after start
		event.Trigger(event.AfterHttpStart, nil)
	}()

	server.shutdown(srv)

	return nil
}

// shutdown
func (server *Server) shutdown(srv *http.Server) {
	// wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	event.Trigger(event.BeforeHttpShutdown, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

