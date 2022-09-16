package http

import (
	"context"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/runtime"
	server2 "github.com/ebar-go/ego/server"
	"github.com/ebar-go/ego/server/protocol"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"sync"
	"time"
)

// HTTPServer represents an HTTP server.
type HTTPServer struct {
	schema protocol.Schema

	instance *http.Server
	initOnce sync.Once

	router                    *gin.Engine
	closeOnce                 sync.Once
	startHooks, shutdownHooks []func()
}

// initialize init http server only once.
func (server *HTTPServer) initialize() {
	server.initOnce.Do(func() {
		server.instance = &http.Server{
			Addr:    server.schema.Bind,
			Handler: server.router,
		}
	})

}

// Serve starts the server.
func (server *HTTPServer) Serve(stop <-chan struct{}) {
	component.Provider().Logger().Infof("listening and serving HTTP on %s", server.schema.Bind)

	for _, hook := range server.startHooks {
		hook()
	}

	go func() {
		if err := server.instance.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			component.Provider().Logger().Fatalf("unable to serve: %v", err)
		}
	}()

	runtime.WaitClose(stop, server.Shutdown)
}

// RegisterRouteLoader registers a route loader
func (server *HTTPServer) RegisterRouteLoader(loader func(router *gin.Engine)) *HTTPServer {
	loader(server.router)
	return server
}

// WithNotFoundHandler provide the handler for not found routes and methods
func (server *HTTPServer) WithNotFoundHandler(notFoundHandler ...gin.HandlerFunc) *HTTPServer {
	server.router.NoRoute(notFoundHandler...)
	server.router.NoMethod(notFoundHandler...)
	return server
}

// EnableCorsMiddleware enables cors middleware
func (server *HTTPServer) EnableCorsMiddleware() *HTTPServer {
	server.router.Use(server2.CORS)
	return server
}

// WithDefaultRecoverMiddleware enables default recover middleware
func (server *HTTPServer) WithDefaultRecoverMiddleware() *HTTPServer {
	server.router.Use(server2.Recover())
	return server
}

// WithDefaultRequestLogMiddleware enables the default request log middleware.
func (server *HTTPServer) WithDefaultRequestLogMiddleware() *HTTPServer {
	server.router.Use(server2.RequestLog())
	return server
}

// EnableTraceMiddleware enables trace middleware with trace header name
func (server *HTTPServer) EnableTraceMiddleware(traceHeader string) *HTTPServer {
	server.router.Use(server2.Trace(traceHeader))
	return server
}

// EnableAvailableHealthCheck enables the health check
func (server *HTTPServer) EnableAvailableHealthCheck() *HTTPServer {
	server.router.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
	return server
}

// EnablePprofHandler enables the profiler for the http server
func (server *HTTPServer) EnablePprofHandler() *HTTPServer {
	pprof.Register(server.router)
	return server
}

// EnableReleaseMode enables the release mode for the http server,it will hide the route tables
func (server *HTTPServer) EnableReleaseMode() *HTTPServer {
	gin.SetMode(gin.ReleaseMode)
	return server
}

// EnableSwaggerHandler enables the swagger handler for the http server
func (server *HTTPServer) EnableSwaggerHandler() *HTTPServer {
	server.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return server
}

// AddStartHook adds a hook function what's called before server is start
func (server *HTTPServer) AddStartHook(hook func()) *HTTPServer {
	server.startHooks = append(server.startHooks, hook)
	return server
}

// AddShutdownHook adds a callback function what's called before the server is shutdown
func (server *HTTPServer) AddShutdownHook(hook func()) *HTTPServer {
	server.shutdownHooks = append(server.shutdownHooks, hook)
	return server
}

// Shutdown shuts down the server.
func (server *HTTPServer) Shutdown() {
	for _, hook := range server.shutdownHooks {
		hook()
	}
	server.closeOnce.Do(server.shutdown)
}

// Shutdown 平滑重启
func (server *HTTPServer) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// stop the server gracefully
	if err := server.instance.Shutdown(ctx); err != nil {
		component.Provider().Logger().Fatalf("HTTPServer shutdown failed: %v", err)
	}
	component.Provider().Logger().Info("HTTPServer showdown success")
}

// NewServer returns a new instance of the HTTPServer.
func NewServer(addr string) *HTTPServer {
	instance := &HTTPServer{
		schema: protocol.NewHttpSchema(addr),
		router: gin.Default(),
	}

	instance.initialize()

	return instance
}
