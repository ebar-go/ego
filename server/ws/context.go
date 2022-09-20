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

func NewContext(conn Conn, body []byte) *Context {
	return &Context{conn: conn, body: body}
}
