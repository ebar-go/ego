package server

import (
	"github.com/ebar-go/ego/rebuild/component"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type RPCServer struct {
	schema Schema

	instance  *grpc.Server
	closeOnce sync.Once
}

// RegisterService registers grpc service
func (server *RPCServer) RegisterService(register func(s *grpc.Server)) *RPCServer {
	register(server.instance)
	return server
}

// Serve start grpc listener
func (server *RPCServer) Serve(stop <-chan struct{}) {
	component.Provider().Logger().Infof("listening and serving GRPC on %s", server.schema.Bind)

	lis, err := net.Listen("tcp", server.schema.Bind)
	if err != nil {
		component.Provider().Logger().Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err := server.instance.Serve(lis); err != nil {
			component.Provider().Logger().Fatalf("failed to serve: %v", err)
		}
	}()

	for {
		select {
		case <-stop:
			server.Shutdown()
		default:

		}
	}
}

func (server *RPCServer) shutdown() {
	server.instance.GracefulStop()
	component.Provider().Logger().Info("RPCServer shutdown success")
}
func (server *RPCServer) Shutdown() {
	server.closeOnce.Do(server.shutdown)
}

func NewGRPCServer(bind string) *RPCServer {
	return &RPCServer{
		schema: Schema{
			Protocol: "grpc",
			Bind:     bind,
		},
		instance: grpc.NewServer(),
	}
}
