package ws

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	server := NewServer(":8082")
	assert.NotNil(t, server)
}

func serveServer(server *Server) {
	ctx, cancel := context.WithCancel(context.Background())
	go server.Serve(ctx.Done())

	time.Sleep(time.Second)
	cancel()
}

func TestServer_Serve(t *testing.T) {
	server := NewServer(":8082")
	serveServer(server)
}
