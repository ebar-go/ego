package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/satori/go.uuid"
)

type IClient interface {
	OnOpen()
	OnClose()
	Read()
	Write()
	SendMessage(message []byte)
}

// Client is a websocket client
type Client struct {
	// 唯一标示
	ID     string

	// 句柄
	Socket *websocket.Conn

	// 发送字符内容
	Send   chan []byte

	// 处理方法
	Handler IHandler

	// 连接时的回调
	OpenHandler func()

	// 关闭时的回调
	CloseHandler func()
	Extends map[string]interface{}
}

var u = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // use default options

// GetUpgradeConnection get web socket connection
func GetUpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	respHeader := http.Header{"Sec-WebSocket-Protocol" :[]string{r.Header.Get("Sec-WebSocket-Protocol")} }
	return  u.Upgrade(w, r, respHeader)
}

// DefaultClient
func DefaultClient(conn *websocket.Conn, handler IHandler) *Client {
	if handler == nil {
		handler = DefaultHandler
	}
	return &Client{ID: uuid.NewV4().String(), Socket: conn, Send: make(chan []byte), Handler:handler, Extends: map[string]interface{}{}}
}

// Send
func (c *Client) SendMessage(message []byte) {
	c.Send <- message
}

// OnConnect
func (c *Client) OnOpen() {
	if c.OpenHandler != nil {
		c.OpenHandler()
	}
}

// OnClose
func (c *Client) OnClose() {
	if c.CloseHandler != nil {
		c.CloseHandler()
	}
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




