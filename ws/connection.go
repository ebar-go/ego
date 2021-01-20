package ws

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// Connection
type Connection struct {
	// unique id
	ID string
	// socket connection
	sockConn *websocket.Conn
	// process handler
	handler Handler
}


var u = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // use default options
// WebsocketConn return web socket connection
func NewConnection(w http.ResponseWriter, r *http.Request) (*Connection, error) {
	respHeader := http.Header{"Sec-WebSocket-Protocol": []string{r.Header.Get("Sec-WebSocket-Protocol")}}
	conn, err := u.Upgrade(w, r, respHeader)
	if err != nil {
		return nil, err
	}

	return &Connection{
		ID:   uuid.NewV4().String(),
		sockConn: conn,
	}, nil
}

// Send 发送数据
func (conn *Connection) Send(message []byte) error{
	return conn.sockConn.WriteMessage(websocket.TextMessage, message)
}

func (conn *Connection) Handle(handler Handler) {
	conn.handler = handler
}

// close
func (c *Connection) close(unregister chan *Connection) {
	_ = c.sockConn.Close()
	unregister <- c
}

// Listen listen connection
func (c *Connection) listen(unregister chan *Connection) {
	defer func() {
		c.close(unregister)
	}()

	for {
		_, message, err := c.sockConn.ReadMessage()
		if err != nil {
			break
		}

		c.handler(message)

	}
}
