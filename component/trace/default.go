package trace

import "github.com/ebar-go/egu"

const (
	prefix = "trace:"
)

var (
	instance = NewTrace()
)

func Id() string {
	return prefix + egu.UUID()
}

func Set(uuid string) {
	instance.Set(uuid)
}

func Get() string {
	return instance.Get()
}

func GC() {
	instance.GC()
}

func Go(fn func()) {
	instance.Go(fn)
}
