package ws

// Context 采用上下文的方式传递
type Context struct {
	extends map[string]interface{}
	message string
}

// Set
func (ctx *Context) Set(key string, value interface{}) {
	if ctx.extends == nil {
		ctx.extends = make(map[string]interface{})
	}
	ctx.extends[key] = value
}

// Get
func (ctx *Context) Get(key string) interface{} {
	if ctx.extends == nil {
		return nil
	}

	return ctx.extends[key]
}

// GetMessage
func (ctx *Context) GetMessage() string {
	return ctx.message
}
