package http

import (
	"context"
	"crypto/tls"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/server/schema"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"sync"
	"time"
)

// Server represents an HTTP server.
type Server struct {
	schema schema.Schema

	instance *http.Server
	initOnce sync.Once

	router                    *gin.Engine
	closeOnce                 sync.Once
	startHooks, shutdownHooks []func()
	tlsConfig                 *tls.Config
}

// Run starts the server.
func (server *Server) Run(stop <-chan struct{}) {
	component.Logger().Infof("listening and serving HTTP on %s", server.schema.Bind)

	for _, hook := range server.startHooks {
		hook()
	}

	go func() {
		if err := server.getInstance().ListenAndServe(); err != nil && err != http.ErrServerClosed {
			component.Logger().Fatalf("unable to serve: %v", err)
		}
	}()

	runtime.WaitClose(stop, server.Shutdown)
}

// RegisterRouteLoader registers a route loader
func (server *Server) RegisterRouteLoader(loader func(router *gin.Engine)) *Server {
	loader(server.router)
	return server
}

// WithNotFoundHandler provide the handler for not found routes and methods
func (server *Server) WithNotFoundHandler(notFoundHandler ...gin.HandlerFunc) *Server {
	server.router.NoRoute(notFoundHandler...)
	server.router.NoMethod(notFoundHandler...)
	return server
}

// EnableCorsMiddleware enables cors middleware
func (server *Server) EnableCorsMiddleware() *Server {
	server.router.Use(CORS)
	return server
}

// WithDefaultRecoverMiddleware enables default recover middleware
func (server *Server) WithDefaultRecoverMiddleware() *Server {
	server.router.Use(Recover())
	return server
}

// WithDefaultRequestLogMiddleware enables the default request log middleware.
func (server *Server) WithDefaultRequestLogMiddleware() *Server {
	server.router.Use(RequestLog())
	return server
}

// EnableTraceMiddleware enables trace middleware with trace header name
func (server *Server) EnableTraceMiddleware(traceHeader string) *Server {
	server.router.Use(Trace(traceHeader))
	return server
}

func (server *Server) EnableCompress() *Server {
	server.router.Use(gzip.Gzip(gzip.DefaultCompression))
	return server
}

// EnableAvailableHealthCheck enables the health check
func (server *Server) EnableAvailableHealthCheck() *Server {
	server.router.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
	return server
}

// EnablePprofHandler enables the profiler for the http server
func (server *Server) EnablePprofHandler() *Server {
	pprof.Register(server.router)
	return server
}

// EnableReleaseMode enables the release mode for the http server,it will hide the route tables
func (server *Server) EnableReleaseMode() *Server {
	gin.SetMode(gin.ReleaseMode)
	return server
}

// EnableSwaggerHandler enables the swagger handler for the http server
func (server *Server) EnableSwaggerHandler() *Server {
	server.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return server
}

// AddStartHook adds a hook function what's called before server is start
func (server *Server) AddStartHook(hook func()) *Server {
	server.startHooks = append(server.startHooks, hook)
	return server
}

// AddShutdownHook adds a callback function what's called before the server is shutdown
func (server *Server) AddShutdownHook(hook func()) *Server {
	server.shutdownHooks = append(server.shutdownHooks, hook)
	return server
}

// EnableTracing enables tracing of jaeger
func (server *Server) EnableTracing() *Server {

	return server
}

// WithTLSConfig sets the tls configuration
func (server *Server) WithTLSConfig(cfg *tls.Config) *Server {
	server.tlsConfig = cfg
	return server
}

// Shutdown shuts down the server.
func (server *Server) Shutdown() {
	for _, hook := range server.shutdownHooks {
		hook()
	}
	server.closeOnce.Do(server.shutdown)
}

// =======================private methods =========================
// getInstance returns the singleton instance of the http.Server.
func (server *Server) getInstance() *http.Server {
	server.initOnce.Do(func() {
		server.instance = &http.Server{
			Addr:      server.schema.Bind,
			Handler:   server.router,
			TLSConfig: server.tlsConfig,
		}
	})
	return server.instance

}

// Shutdown 平滑重启
func (server *Server) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// stop the server gracefully
	if err := server.getInstance().Shutdown(ctx); err != nil {
		component.Logger().Fatalf("Server shutdown failed: %v", err)
	} else {
		component.Logger().Info("Server showdown success")
	}

}

// NewServer returns a new instance of the Server.
func NewServer(addr string) *Server {
	instance := &Server{
		schema: schema.NewHttpSchema(addr),
		router: gin.Default(),
	}

	return instance
}
