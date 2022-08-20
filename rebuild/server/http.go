package server

import (
	"context"
	"github.com/ebar-go/ego/rebuild/component"
	"github.com/ebar-go/ego/rebuild/runtime"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
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
