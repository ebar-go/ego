package server

import (
	"context"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/runtime"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"sync"
	"time"
)

type HTTPServer struct {
	schema Schema

	instance *http.Server

	router    *gin.Engine
	closeOnce sync.Once
}

func (server *HTTPServer) Serve(stop <-chan struct{}) {
	component.Provider().Logger().Infof("listening and serving HTTP on %s", server.schema.Bind)

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
	server.router.Use(CORS)
	return server
}

// WithDefaultRecoverMiddleware enables default recover middleware
func (server *HTTPServer) WithDefaultRecoverMiddleware() *HTTPServer {
	server.router.Use(Recover())
	return server
}

// EnableTraceMiddleware enables trace middleware with trace header name
func (server *HTTPServer) EnableTraceMiddleware(traceHeader string) *HTTPServer {
	server.router.Use(Trace(traceHeader))
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

// Shutdown shuts down the server.
func (server *HTTPServer) Shutdown() {
	server.closeOnce.Do(server.shutdown)
}

// Shutdown 平滑重启
func (server *HTTPServer) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// stop the server gracefully
	if err := server.instance.Shutdown(ctx); err != nil {
		component.Provider().Logger().Fatalf("HTTPServer shutdown failed:", err)
	}
	component.Provider().Logger().Info("HTTPServer showdown success")
}

func NewHTTPServer(addr string) *HTTPServer {
	router := gin.Default()
	return &HTTPServer{
		schema: Schema{
			Protocol: "http",
			Bind:     addr,
		},
		router: router,
		instance: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}
