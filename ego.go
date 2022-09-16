package ego

import (
	"github.com/ebar-go/ego/server"
)

// New creates a new NamedEngine instance with default name.
func New() *NamedEngine {
	return NewNamedEngine("default")
}

// NewHttpServer creates a new http server instance.
func NewHTTPServer(addr string) *server.HTTPServer {
	return server.NewHTTPServer(addr)
}

// NewGrpcServer creates a new grpc server instance.
func NewGRPCServer(addr string) *server.RPCServer {
	return server.NewGRPCServer(addr)
}

// NewWebSocketServer creates a new web server instance.
func NewWebsocketServer(addr string) *server.WebSocketServer {
	return server.NewWebSocketServer(addr)
}
