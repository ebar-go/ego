package grpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	server := NewServer(":8081")
	assert.NotNil(t, server)
}

func serveServer(server *Server) {
	ctx, cancel := context.WithCancel(context.Background())
	go server.Run(ctx.Done())

	time.Sleep(time.Second)
	cancel()
}
func TestServer_Serve(t *testing.T) {
	server := NewServer(":8081")
	serveServer(server)
}
