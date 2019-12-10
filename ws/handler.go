package ws

type IHandler func(ctx *Context) string

// Context 采用上下文的方式传递
type Context struct {
	extends map[string]interface{}
	message string
}

func (ctx *Context) Set(key string, value interface{}) {
	if ctx.extends == nil {
		ctx.extends = make(map[string]interface{})
	}
	ctx.extends[key] = value
}

func (ctx *Context) Get(key string) interface{} {
	if ctx.extends == nil {
		return nil
	}

	return ctx.extends[key]
}

func (ctx *Context) GetMessage() string {
	return ctx.message
}


func DefaultHandler(ctx *Context) string {
	// do nothing
	return ctx.message
}