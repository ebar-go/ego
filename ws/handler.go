package ws

type IHandler func(ctx *Context) string

// DefaultHandler
func DefaultHandler(ctx *Context) string {
	// do nothing
	return ctx.message
}
