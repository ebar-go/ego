package ws

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	uuid "github.com/satori/go.uuid"
	"net"
)

// Conn websocket connection
type Conn interface {
	ID() string
	Push(msg []byte) error
	IP() string
}

func NewWrapConnection(netConn net.Conn) *connection {
	return &connection{conn: netConn, uuid: uuid.NewV4().String()}
}

type connection struct {
	conn net.Conn
	uuid string
}

func (c *connection) ID() string {
	return c.uuid
}

// Push send message to client
func (c *connection) Push(msg []byte) error {
	return wsutil.WriteServerBinary(c.conn, msg)
}

func (c *connection) IP() string {
	return c.conn.RemoteAddr().String()
}

func (c *connection) readClientData() ([]byte, error) {

	msg, op, err := wsutil.ReadClientData(c.conn)
	if err != nil {
		return nil, err
	}
	if op != ws.OpBinary && op != ws.OpText {
		return nil, nil
	}

	return msg, err
}