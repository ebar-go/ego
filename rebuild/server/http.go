package server

import (
	"context"
	"github.com/ebar-go/ego/rebuild/component"
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
	srv := &http.Server{
		Addr:    server.schema.Bind,
		Handler: server.router,
	}

	component.Provider().Logger().Infof("listening and serving HTTP on %s\n", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			component.Provider().Logger().Fatalf("unable to serve: %v\n", err)
		}
	}()

	server.instance = srv
	<-stop
	server.Shutdown()
}

// Shutdown 平滑重启
func (server *HTTPServer) shutdown() {
	if server.instance == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.instance.Shutdown(ctx); err != nil {
		component.Provider().Logger().Fatalf("HTTPServer shutdown failed:", err)
	}
	component.Provider().Logger().Info("HTTPServer showdown success")
}

func (server *HTTPServer) Shutdown() {
	server.closeOnce.Do(server.shutdown)
}

func NewHTTPServer(addr string) *HTTPServer {
	return &HTTPServer{
		schema: Schema{
			Protocol: "http",
			Bind:     addr,
		},
		router: gin.Default(),
	}
}
