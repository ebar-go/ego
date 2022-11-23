package ego

import (
	"github.com/ebar-go/ego/protocol/grpc"
	"github.com/ebar-go/ego/protocol/http"
)

// New creates a new Aggregator instance with default name.
func New() *Aggregator {
	return NewAggregator("default")
}

// NewHttpServer creates a new http server instance.
func NewHTTPServer(addr string) *http.Server {
	return http.NewServer(addr)
}

// NewGRPCServer creates a new grpc server instance.
func NewGRPCServer(addr string) *grpc.Server {
	return grpc.NewServer(addr)
}
