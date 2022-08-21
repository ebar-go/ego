package ego

import (
	"github.com/ebar-go/ego/aggregator"
	"github.com/ebar-go/ego/server"
)

// NewAggregator creates a new Aggregator instance
func NewAggregatorServer() *aggregator.Aggregator {
	return aggregator.NewAggregator()
}

func NewHTTPServer(addr string) *server.HTTPServer {
	return server.NewHTTPServer(addr)
}

func NewGRPCServer(addr string) *server.RPCServer {
	return server.NewGRPCServer(addr)
}

func NewWebsocketServer(addr string) *server.WebSocketServer {
	return server.NewWebSocketServer(addr)
}
