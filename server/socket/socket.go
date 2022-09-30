package socket

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/server/protocol"
	"github.com/gobwas/ws"
	"log"
)

func NewWebsocket(bind string) *Server {
	return &Server{
		schema: protocol.NewWSSchema(bind),
		upgrader: ws.Upgrader{
			OnHeader: func(key, value []byte) (err error) {
				log.Printf("non-websocket header: %q=%q", key, value)
				return
			},
		},
		connectHandler:    func(conn Connection) {},
		disconnectHandler: func(conn Connection) {},
		requestHandler: func(ctx *Context) {
			component.Provider().Logger().Infof("received request: %v", string(ctx.Body()))
		},
		worker: 1,
	}
}

func NewTCP(bind string) *Server {
	return &Server{
		schema:            protocol.NewTCPSchema(bind),
		connectHandler:    func(conn Connection) {},
		disconnectHandler: func(conn Connection) {},
		requestHandler: func(ctx *Context) {
			component.Provider().Logger().Infof("received request: %v", string(ctx.Body()))
		},
		worker: 1,
	}
}
