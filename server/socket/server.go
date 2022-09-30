package socket

import (
	"context"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server/protocol"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/gobwas/ws"
	"net"
	"sync"
)

type Server struct {
	connectHandler    func(conn Connection)
	disconnectHandler func(conn Connection)
	requestHandler    func(ctx *Context)
	worker            int

	schema    protocol.Schema
	closeOnce sync.Once
	stopCh    chan struct{}

	upgrader ws.Upgrader
}

// OnConnect is set callback when the connection is established
func (server *Server) OnConnect(handler func(conn Connection)) *Server {
	server.connectHandler = handler
	return server
}

// handleConnect is called when the connection is established
func (server *Server) handleConnect(conn Connection) {
	if server.connectHandler == nil {
		return
	}
	server.connectHandler(conn)
}

// OnMessage is set callback  when a message is received.
func (server *Server) OnMessage(handler func(ctx *Context)) *Server {
	if handler != nil {
		server.requestHandler = handler
	}
	return server
}

// OnDisconnect is set callback when the client disconnects from the server
func (server *Server) OnDisconnect(handler func(conn Connection)) *Server {
	server.disconnectHandler = handler
	return server
}

// WithWorker sets the concurrent worker threads.
func (server *Server) WithWorker(worker int) *Server {
	server.worker = worker
	return server
}

// handleConnect is called when the client disconnects from the server
func (server *Server) handleDisconnect(conn Connection) {
	if server.disconnectHandler == nil {
		return
	}
	server.disconnectHandler(conn)
}

func (server *Server) shutdown() {
	close(server.stopCh)
	component.Provider().Logger().Info("Server shutdown success")
}

// Shutdown stop websocket server.
func (server *Server) Shutdown() {
	server.closeOnce.Do(server.shutdown)
}

// Serve start websocket server
func (server *Server) Serve(stop <-chan struct{}) {
	component.Provider().Logger().Infof("listening and serving websocket on %s", server.schema.Bind)

	listener, err := net.Listen("tcp", server.schema.Bind)
	if err != nil {
		component.Provider().Logger().Fatalf("failed to listen websocket: %v", err)
	}

	// cancel function is used to be called when the server need shutdown
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-server.stopCh
		cancel()
	}()
	for i := 0; i < server.worker; i++ {
		go server.serve(listener, ctx.Done())
	}

	runtime.WaitClose(stop, server.Shutdown)
}

func (server *Server) serve(listener net.Listener, stopCh <-chan struct{}) {
	var err error
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		if server.schema.Protocol == protocol.WS {
			err = server.acceptWebsocket(listener)
		} else if server.schema.Protocol == protocol.TCP {
			err = server.acceptTCP(listener)
		}
		if err != nil {
			component.Provider().Logger().Errorf("failed to accept: %v", err)
		}
	}
}

func (server *Server) acceptWebsocket(listener net.Listener) error {
	conn, err := listener.Accept()
	if err != nil {
		return errors.WithMessage(err, "listener.Accept")
	}
	_, err = server.upgrader.Upgrade(conn)
	if err != nil {
		return errors.WithMessage(err, "listener.Upgrade")
	}

	connection := NewWrapConnection(conn)
	server.handleConnect(connection)

	go func() {
		defer func() {
			server.handleDisconnect(connection)
			_ = conn.Close()
		}()

		for {
			msg, err := connection.readClientData()
			if err != nil {
				return
			}

			if len(msg) == 0 {
				continue
			}

			ctx := NewContext(connection, msg)
			server.requestHandler(ctx)
		}

	}()
	return nil
}

func (server *Server) acceptTCP(listener net.Listener) error {
	return nil
}
