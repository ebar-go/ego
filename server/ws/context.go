package ws

import "context"

type Context struct {
	context.Context
	body []byte
	conn Conn
}

func (ctx *Context) Body() []byte {
	return ctx.body
}

func (ctx *Context) Output(msg []byte) {
	_ = ctx.conn.Push(msg)
}

func (ctx *Context) Conn() Conn {
	return ctx.conn
}

func NewContext(conn Conn, body []byte) *Context {
	return &Context{Context: context.Background(), conn: conn, body: body}
}
