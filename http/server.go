package http

import (
	"context"
	"fmt"
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/event"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/ego/http/validator"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// Server Web服务管理器
type Server struct {
	// gin engine
	router *gin.Engine
	// server config
	conf *config.Config
}

// HttpServer 获取Server示例
func New(conf *config.Config) *Server {
	router := gin.New()

	// use global trace middleware
	router.Use(middleware.Trace(conf.TraceHeader))
	router.Use(func(ctx *gin.Context) {
		if ctx.Request.RequestURI == "/favicon.ico" {
			ctx.AbortWithStatus(200)
			return
		}
		ctx.Next()
	})

	// 404
	router.NoRoute(notFoundHandler)
	router.NoMethod(notFoundHandler)

	return &Server{
		conf: conf,
		router: router,
	}
}

func (server *Server) Router() *gin.Engine {
	return server.router
}

func (server *Server) beforeStart()  {
	binding.Validator = new(validator.Validator)
	// before start
	event.Trigger(event.BeforeHttpStart, nil)

	if server.conf.Pprof {
		pprof.Register(server.router)
	}

	if server.conf.Swagger {
		ginSwagger.WrapHandler(swaggerFiles.Handler)
	}
}

// Serve run http server
func (server *Server) Serve() error {
	server.beforeStart()

	completeHost := net.JoinHostPort("", strconv.Itoa(server.conf.Port))

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

	server.listen(srv)

	return nil
}

// shutdown
func (server *Server) listen(srv *http.Server) {
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


// notFoundHandler 404
func notFoundHandler(ctx *gin.Context) {
	response.WrapContext(ctx).Error(404,
		fmt.Sprintf("404 Not Found: %s", ctx.Request.RequestURI))
}

