package ws

import (
	"github.com/gobwas/ws/wsutil"
	"net"
)

// Conn websocket connection
type Conn struct {
	conn net.Conn
}

// Push send message to client
func (c *Conn) Push(msg []byte) error {
	return wsutil.WriteServerBinary(c.conn, msg)
}

func (c *Conn) IP() string {
	return c.conn.RemoteAddr().String()
}
