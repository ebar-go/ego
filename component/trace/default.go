package trace

import (
	uuid "github.com/satori/go.uuid"
)

const (
	prefix = "trace:"
)

var (
	instance = NewTrace()
)

func Id() string {
	return prefix + uuid.NewV4().String()
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
