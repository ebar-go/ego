package server

import (
	"context"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/runtime"
	"github.com/gobwas/ws"
	"log"
	"net"
	"sync"
)

type WebSocketServer struct {
	schema    Schema
	closeOnce sync.Once
	upgrader  ws.Upgrader
	listener  net.Listener
	cancel    context.CancelFunc
}

func (server *WebSocketServer) Serve(stop <-chan struct{}) {
	component.Provider().Logger().Infof("listening and serving websocket on %s", server.schema.Bind)

	ln, err := net.Listen("tcp", server.schema.Bind)
	if err != nil {
		component.Provider().Logger().Fatalf("failed to listen: %v", err)
	}
	server.listener = ln

	ctx, cancel := context.WithCancel(context.Background())
	server.cancel = cancel
	go func() {
		for {
			select {
			case <-ctx.Done():
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

func (server *WebSocketServer) Shutdown() {
	server.closeOnce.Do(server.shutdown)
}

func (server *WebSocketServer) accept() error {
	conn, err := server.listener.Accept()
	if err != nil {
		return errors.WithMessage(err, "listener.Accept")
	}

	_, err = server.upgrader.Upgrade(conn)
	if err != nil {
		return errors.WithMessage(err, "listener.Upgrade")
	}
	return nil
}

func NewWebSocketServer(bind string) *WebSocketServer {
	return &WebSocketServer{
		schema: Schema{Protocol: "ws", Bind: bind},
		upgrader: ws.Upgrader{
			OnHeader: func(key, value []byte) (err error) {
				log.Printf("non-websocket header: %q=%q", key, value)
				return
			},
		},
	}
}
