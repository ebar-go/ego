package server

import (
	"context"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/runtime"
	"github.com/ebar-go/ego/server/protocol"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"sync"
)

// Conn websocket connection
type Conn struct {
	conn net.Conn
}

// Push send message to client
func (c *Conn) Push(msg []byte) error {
	return wsutil.WriteServerBinary(c.conn, msg)
}

// WebSocketServer provide the websocket server.
type WebSocketServer struct {
	schema            protocol.Schema
	closeOnce         sync.Once
	upgrader          ws.Upgrader
	listener          net.Listener
	cancel            context.CancelFunc
	connectHandler    func(conn Conn)
	disconnectHandler func(conn Conn)
	requestHandler    func(conn Conn, msg []byte)
}

// OnConnect is set callback when the connection is established
func (server *WebSocketServer) OnConnect(handler func(conn Conn)) *WebSocketServer {
	server.connectHandler = handler
	return server
}

// handleConnect is called when the connection is established
func (server *WebSocketServer) handleConnect(conn Conn) {
	if server.connectHandler == nil {
		return
	}
	server.connectHandler(conn)
}

// OnMessage is set callback  when a message is received.
func (server *WebSocketServer) OnMessage(handler func(conn Conn, msg []byte)) *WebSocketServer {
	if handler != nil {
		server.requestHandler = handler
	}
	return server
}

// OnDisconnect is set callback when the client disconnects from the server
func (server *WebSocketServer) OnDisconnect(handler func(conn Conn)) *WebSocketServer {
	server.disconnectHandler = handler
	return server
}

// handleConnect is called when the client disconnects from the server
func (server *WebSocketServer) handleDisconnect(conn Conn) {
	if server.disconnectHandler == nil {
		return
	}
	server.disconnectHandler(conn)
}

// Serve start websocket server
func (server *WebSocketServer) Serve(stop <-chan struct{}) {
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

func (server *WebSocketServer) shutdown() {
	server.cancel()
	component.Provider().Logger().Info("WebSocketServer shutdown success")
}

// Shutdown stop websocket server.
func (server *WebSocketServer) Shutdown() {
	server.closeOnce.Do(server.shutdown)
}

// accept process incoming connection, and trigger callback with OnConnect,OnDisconnect,OnMessage.
func (server *WebSocketServer) accept() error {
	conn, err := server.listener.Accept()
	if err != nil {
		return errors.WithMessage(err, "listener.Accept")
	}

	_, err = server.upgrader.Upgrade(conn)
	if err != nil {
		return errors.WithMessage(err, "listener.Upgrade")
	}

	connection := Conn{conn: conn}
	server.handleConnect(connection)

	go func() {
		defer func() {
			server.handleDisconnect(connection)
			_ = conn.Close()
		}()

		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				return
			}

			if op != ws.OpBinary {
				continue
			}

			server.requestHandler(connection, msg)
		}

	}()
	return nil
}

func NewWebSocketServer(bind string) *WebSocketServer {
	return &WebSocketServer{
		schema: protocol.Schema{Protocol: "ws", Bind: bind},
		upgrader: ws.Upgrader{
			OnHeader: func(key, value []byte) (err error) {
				log.Printf("non-websocket header: %q=%q", key, value)
				return
			},
		},
		requestHandler: func(conn Conn, msg []byte) {
			component.Provider().Logger().Infof("received request: %v", string(msg))
		},
	}
}
