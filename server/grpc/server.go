package grpc

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/runtime"
	"github.com/ebar-go/ego/server/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"sync"
)

// Server represents a gRPC server.
type Server struct {
	schema protocol.Schema

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

// Serve start grpc listener
func (server *Server) Serve(stop <-chan struct{}) {
	component.Provider().Logger().Infof("listening and serving GRPC on %s", server.schema.Bind)

	lis, err := net.Listen("tcp", server.schema.Bind)
	if err != nil {
		component.Provider().Logger().Fatalf("failed to listen rpc: %v", err)
	}

	go func() {
		if err := server.getInstance().Serve(lis); err != nil {
			component.Provider().Logger().Fatalf("failed to serve: %v", err)
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
	component.Provider().Logger().Info("Server shutdown success")
}

// NewServer returns a new instance of the Server.
func NewServer(bind string) *Server {
	return &Server{
		schema: protocol.NewGRPCSchema(bind),
	}
}