package http

import (
	"context"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	server := NewServer(":8081").
		RegisterRouteLoader(func(router *gin.Engine) {
			router.Any("test", func(c *gin.Context) {
				c.JSON(200, gin.H{"name": "123"})
			})
		})
	ctx, cancel := context.WithCancel(context.Background())
	go server.Run(ctx.Done())
	assert.NotNil(t, server)
	runtime.Shutdown(func() {
		cancel()
	})
}

func serveServer(server *Server) {
	ctx, cancel := context.WithCancel(context.Background())
	go server.Run(ctx.Done())

	time.Sleep(time.Second)
	cancel()
}

func TestHTTPServer_Serve(t *testing.T) {
	server := NewServer(":8080")
	serveServer(server)
}

func TestHTTPServer_AddStartHook(t *testing.T) {
	server := NewServer(":8080").AddStartHook(func() {
		log.Println("hook before start")
	})

	serveServer(server)
}

func TestHTTPServer_AddShutdownHook(t *testing.T) {
	server := NewServer(":8080").AddShutdownHook(func() {
		log.Println("hook before shutdown")
	})

	serveServer(server)
}
