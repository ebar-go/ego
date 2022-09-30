package ego

import (
	"github.com/ebar-go/ego/server/grpc"
	"github.com/ebar-go/ego/server/http"
	"github.com/ebar-go/ego/server/socket"
)

// New creates a new NamedEngine instance with default name.
func New() *NamedEngine {
	return NewNamedEngine("default")
}

// NewHttpServer creates a new http server instance.
func NewHTTPServer(addr string) *http.Server {
	return http.NewServer(addr)
}

// NewGRPCServer creates a new grpc server instance.
func NewGRPCServer(addr string) *grpc.Server {
	return grpc.NewServer(addr)
}

// NewWebSocketServer creates a new web server instance.
func NewWebsocketServer(addr string) *socket.Server {
	return socket.NewWebsocket(addr)
}
