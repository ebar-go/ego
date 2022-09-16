package server

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/runtime"
	"github.com/ebar-go/ego/server/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"sync"
)

type RPCServer struct {
	schema protocol.Schema

	initOnce  sync.Once
	instance  *grpc.Server
	closeOnce sync.Once
	options   []grpc.ServerOption
}

// WithOptions sets the options for the RPC server.It must be called before RegisterService.
func (server *RPCServer) WithOptions(opts ...grpc.ServerOption) *RPCServer {
	server.options = append(server.options, opts...)
	return server
}

// WithKeepAliveParams sets the KeepAlive option.It must be called before RegisterService.
func (server *RPCServer) WithKeepAliveParams(kp keepalive.ServerParameters) *RPCServer {
	return server.WithOptions(grpc.KeepaliveParams(kp))
}

// WithChainUnaryInterceptor sets the interceptors.It must be called before RegisterService.
func (server *RPCServer) WithChainUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) *RPCServer {
	return server.WithOptions(grpc.ChainUnaryInterceptor())
}

// RegisterService registers grpc service
func (server *RPCServer) RegisterService(register func(s *grpc.Server)) *RPCServer {
	register(server.getInstance())
	return server
}

// getInstance returns the singleton instance of the grpc server
func (server *RPCServer) getInstance() *grpc.Server {
	server.initOnce.Do(func() {
		server.instance = grpc.NewServer(server.options...)
	})
	return server.instance
}

// Serve start grpc listener
func (server *RPCServer) Serve(stop <-chan struct{}) {
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
func (server *RPCServer) Shutdown() {
	server.closeOnce.Do(server.shutdown)
}

func (server *RPCServer) shutdown() {
	// stop grpc server gracefully
	server.instance.GracefulStop()
	component.Provider().Logger().Info("RPCServer shutdown success")
}

func NewGRPCServer(bind string) *RPCServer {
	return &RPCServer{
		schema: protocol.Schema{
			Protocol: "grpc",
			Bind:     bind,
		},
	}
}
