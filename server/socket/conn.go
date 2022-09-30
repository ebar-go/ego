package socket

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	uuid "github.com/satori/go.uuid"
	"net"
	"sync"
)

// Connection define socket connection interface
type Connection interface {
	// ID returns socket connection uuid
	ID() string

	// IP returns client ip address
	IP() string

	// Close closes socket connection
	Close() error

	// Push send message to client.
	Push(msg []byte) error

	// Property returns client property
	Property() *Property
}

type connection struct {
	conn     net.Conn
	uuid     string
	property *Property
}

func (c *connection) Close() error {
	return c.conn.Close()
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

func NewWrapConnection(netConn net.Conn) *connection {
	return &connection{conn: netConn, uuid: uuid.NewV4().String(), property: &Property{properties: map[string]any{}}}
}

type Property struct {
	mu         sync.RWMutex // guards the properties
	properties map[string]any
}

func (p *Property) Set(key string, value any) {
	p.mu.Lock()
	p.properties[key] = value
	p.mu.Unlock()
}

func (p *Property) Get(key string) any {
	p.mu.RLock()
	property := p.properties[key]
	p.mu.RUnlock()
	return property
}

func (p *Property) GetString(key string) string {
	property := p.Get(key)
	if property == nil {
		return ""
	}
	return property.(string)
}
