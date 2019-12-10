package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/satori/go.uuid"
)

type IClient interface {
	Read()
	Write()
}

// Client is a websocket client
type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
	Handler IHandler
	Extends map[string]interface{}
}

var u = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // use default options

// GetUpgradeConnection get web socket connection
func GetUpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return  u.Upgrade(w, r, nil)
}

// DefaultClient
func DefaultClient(conn *websocket.Conn, handler IHandler) *Client {
	if handler == nil {
		handler = DefaultHandler
	}
	return &Client{ID: uuid.NewV4().String(), Socket: conn, Send: make(chan []byte), Handler:handler, Extends: map[string]interface{}{}}
}



// Read 读取数据
func (c *Client) Read() {
	defer func() {
		manager.UnregisterClient(c)
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			manager.UnregisterClient(c)
			c.Socket.Close()
			break
		}

		ctx := &Context{message: string(message)}

		// merge client extends
		for key, value := range c.Extends {
			ctx.Set(key, value)
		}

		// get response and send
		manager.Broadcast(c.Handler(ctx))
	}
}

// Write 写数据
func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}




