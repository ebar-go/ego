package ws

import (
	"context"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server/protocol"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/gobwas/ws"
	"log"
	"net"
	"sync"
)

// Server provide the websocket server.
type Server struct {
	schema            protocol.Schema
	closeOnce         sync.Once
	upgrader          ws.Upgrader
	listener          net.Listener
	cancel            context.CancelFunc
	connectHandler    func(conn Conn)
	disconnectHandler func(conn Conn)
	requestHandler    func(ctx *Context)
}

// OnConnect is set callback when the connection is established
func (server *Server) OnConnect(handler func(conn Conn)) *Server {
	server.connectHandler = handler
	return server
}

// handleConnect is called when the connection is established
func (server *Server) handleConnect(conn Conn) {
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
func (server *Server) OnDisconnect(handler func(conn Conn)) *Server {
	server.disconnectHandler = handler
	return server
}

// handleConnect is called when the client disconnects from the server
func (server *Server) handleDisconnect(conn Conn) {
	if server.disconnectHandler == nil {
		return
	}
	server.disconnectHandler(conn)
}

// Serve start websocket server
func (server *Server) Serve(stop <-chan struct{}) {
	component.Provider().Logger().Infof("listening and serving websocket on %s", server.schema.Bind)

	ln, err := net.Listen("tcp", server.schema.Bind)
	if err != nil {
		component.Provider().Logger().Fatalf("failed to listen websocket: %v", err)
	}
	server.listener = ln

	// cancel function is used to be called when the server need shutdown
	ctx, cancel := context.WithCancel(context.Background())
	server.cancel = cancel
	go func() {
		for {
			select {
			case <-ctx.Done(): // if cancel is called, goroutine should exit.
				return
			default:
				if err := server.accept(); err != nil {
					component.Provider().Logger().Errorf("failed to accept: %v", err)
					continue
				}
			}
		}

	}()

	runtime.WaitClose(stop, server.Shutdown)
}

func (server *Server) shutdown() {
	server.cancel()
	component.Provider().Logger().Info("Server shutdown success")
}

// Shutdown stop websocket server.
func (server *Server) Shutdown() {
	server.closeOnce.Do(server.shutdown)
}

// accept process incoming connection, and trigger callback with OnConnect,OnDisconnect,OnMessage.
func (server *Server) accept() error {
	conn, err := server.listener.Accept()
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

func NewServer(bind string) *Server {
	return &Server{
		schema: protocol.NewWSSchema(bind),
		upgrader: ws.Upgrader{
			OnHeader: func(key, value []byte) (err error) {
				log.Printf("non-websocket header: %q=%q", key, value)
				return
			},
		},
		connectHandler:    func(conn Conn) {},
		disconnectHandler: func(conn Conn) {},
		requestHandler: func(ctx *Context) {
			component.Provider().Logger().Infof("received request: %v", string(ctx.Body()))
		},
	}
}
