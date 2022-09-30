package socket

import "context"

type Context struct {
	context.Context
	body []byte
	conn Connection
}

func (ctx *Context) Body() []byte {
	return ctx.body
}

func (ctx *Context) Output(msg []byte) {
	_ = ctx.conn.Push(msg)
}

func (ctx *Context) Conn() Connection {
	return ctx.conn
}

func NewContext(conn Connection, body []byte) *Context {
	return &Context{Context: context.Background(), conn: conn, body: body}
}
