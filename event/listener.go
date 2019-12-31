package event


// Listener
type Listener struct {
	Handler Handler
}

// 监听器函数
type Handler func(ev Event)

// NewListener
func NewListener(h Handler) *Listener {
	l := new(Listener)
	l.Handler = h
	return l
}
