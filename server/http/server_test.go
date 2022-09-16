package http

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	server := NewServer(":8080")
	assert.NotNil(t, server)
}

func serveServer(server *Server) {
	ctx, cancel := context.WithCancel(context.Background())
	go server.Serve(ctx.Done())

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
