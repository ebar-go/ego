package http

import (
	"context"
	"fmt"
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
	"strconv"
	"time"
)

// Server Web服务管理器
type Server struct {
	// gin engine
	router *gin.Engine
	// server config
	conf *Config
	// httpServer
	instance *http.Server
}

// NewServer 获取Server示例
func NewServer(conf *Config) *Server {
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
		conf:   conf,
		router: router,
	}
}

func (server *Server) registerPprof() {
	if server.conf.Pprof {
		pprof.Register(server.router)
	}
}

func (server *Server) handleSwagger() {
	if server.conf.Swagger {
		server.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func (server *Server) beforeStart() {
	binding.Validator = new(validator.Validator)
	// before start
	event.Trigger(event.BeforeHttpStart, nil)

	// register pprof
	server.registerPprof()
	// handle swagger api
	server.handleSwagger()
}

// Run run http server
func (server *Server) Serve() {
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
			log.Fatalf("%v\n", err)
		}

		// after start
		event.Trigger(event.AfterHttpStart, nil)
	}()

	server.instance = srv
}

// Close 平滑重启
func (server *Server) Close() {
	if server.instance == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.instance.Shutdown(ctx); err != nil {
		log.Fatal("HTTP shutdown:", err)
	}
	log.Println("HTTP showdown")
}

// notFoundHandler 404
func notFoundHandler(ctx *gin.Context) {
	response.WrapContext(ctx).Error(404,
		fmt.Sprintf("404 Not Found: %s", ctx.Request.RequestURI))
}
