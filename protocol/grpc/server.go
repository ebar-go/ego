package grpc

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/protocol/schema"
	"github.com/ebar-go/ego/utils/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"sync"
)

// Server represents a gRPC server.
type Server struct {
	schema schema.Schema

	initOnce  sync.Once
	instance  *grpc.Server
	closeOnce sync.Once
	options   []grpc.ServerOption
}

// WithOptions sets the options for the RPC server.It must be called before RegisterService.
func (server *Server) WithOptions(opts ...grpc.ServerOption) *Server {
	server.options = append(server.options, opts...)
	return server
}

// WithKeepAliveParams sets the KeepAlive option.It must be called before RegisterService.
func (server *Server) WithKeepAliveParams(kp keepalive.ServerParameters) *Server {
	return server.WithOptions(grpc.KeepaliveParams(kp))
}

// WithChainUnaryInterceptor sets the interceptors.It must be called before RegisterService.
func (server *Server) WithChainUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) *Server {
	return server.WithOptions(grpc.ChainUnaryInterceptor())
}

// RegisterService registers grpc service
func (server *Server) RegisterService(register func(s *grpc.Server)) *Server {
	register(server.getInstance())
	return server
}

// Run start grpc listener
func (server *Server) Run(stop <-chan struct{}) {
	component.Logger().Infof("listening and serving GRPC on %s", server.schema.Bind)

	lis, err := net.Listen("tcp", server.schema.Bind)
	if err != nil {
		component.Logger().Fatalf("failed to listen rpc: %v", err)
	}

	go func() {
		if err := server.getInstance().Serve(lis); err != nil {
			component.Logger().Fatalf("failed to serve: %v", err)
		}
	}()

	runtime.WaitClose(stop, server.Shutdown)
}

// Shutdown shuts down the server.
func (server *Server) Shutdown() {
	server.closeOnce.Do(server.shutdown)
}

// ========================= private methods =========================

// getInstance returns the singleton instance of the grpc server
func (server *Server) getInstance() *grpc.Server {
	server.initOnce.Do(func() {
		server.instance = grpc.NewServer(server.options...)
	})
	return server.instance
}

func (server *Server) shutdown() {
	// stop grpc server gracefully
	server.instance.GracefulStop()
	component.Logger().Info("Server shutdown success")
}

// NewServer returns a new instance of the Server.
func NewServer(address string) *Server {
	return &Server{
		schema: schema.NewGRPCSchema(address),
	}
}
