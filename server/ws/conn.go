package ws

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	uuid "github.com/satori/go.uuid"
	"net"
	"sync"
)

// Conn websocket connection
type Conn interface {
	ID() string
	Push(msg []byte) error
	IP() string
	Property() *Property
}

func NewWrapConnection(netConn net.Conn) *connection {
	return &connection{conn: netConn, uuid: uuid.NewV4().String(), property: &Property{properties: map[string]any{}}}
}

type connection struct {
	conn     net.Conn
	uuid     string
	property *Property
}

func (c *connection) Property() *Property {
	return c.property
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

type Property struct {
	mu         sync.Mutex //
	properties map[string]any
}

func (p *Property) Set(key string, value any) {
	p.mu.Lock()
	p.properties[key] = value
	p.mu.Unlock()
}

func (p *Property) Get(key string) any {
	p.mu.Lock()
	property := p.properties[key]
	p.mu.Unlock()
	return property
}

func (p *Property) GetString(key string) string {
	property := p.Get(key)
	if property == nil {
		return ""
	}
	return property.(string)
}
