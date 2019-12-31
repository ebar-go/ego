package ws

import (
	"github.com/ebar-go/ego/helper"
	"github.com/gorilla/websocket"
)

// IClient
type IClient interface {
	OnOpen()
	OnClose()
	Listen()
	SendMessage(message []byte)
}

// Client is a websocket client
type Client struct {
	// unique id
	ID string

	// connection
	Conn *websocket.Conn

	// 处理方法
	Handler IHandler

	// 连接时的回调
	OpenHandler func()

	// 关闭时的回调
	CloseHandler func()

	// 扩展字段
	Extends map[string]interface{}

	manager Manager
}

// NewClient return Client
func NewClient(conn *websocket.Conn, handler IHandler) *Client {
	if handler == nil {
		handler = DefaultHandler
	}
	return &Client{
		ID:      uuid.NewV4().String(),
		Conn:    conn,
		Handler: handler,
		Extends: map[string]interface{}{},
	}
}

// SendMessage
func (c *Client) SendMessage(message []byte) {
	_ = c.Conn.WriteMessage(websocket.TextMessage, message)

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

// 关闭连接
func (c *Client) close() {
	c.OnClose()
	_ = c.Conn.Close()
	c.manager.UnregisterClient(c)
}

// Listen 监听
func (c *Client) Listen() {
	defer func() {
		c.close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			helper.Debug(err)
			break
		}

		ctx := &Context{message: string(message)}

		// merge client extends
		for key, value := range c.Extends {
			ctx.Set(key, value)
		}

		result := c.Handler(ctx)
		c.SendMessage([]byte(result))
	}
}
