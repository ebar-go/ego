package runtime

import (
	"log"
	"runtime"
)

var PanicHandlers = []func(interface{}){logPanic}

var (
	reallyCrash = true
)

func SetReallyCrash(really bool) {
	reallyCrash = really
}

func HandleCrash(additionalHandlers ...func(interface{})) {
	if r := recover(); r != nil {
		for _, fn := range PanicHandlers {
			fn(r)
		}
		for _, fn := range additionalHandlers {
			fn(r)
		}
		if reallyCrash {
			// Actually proceed to panic.
			panic(r)
		}
	}
}

func logPanic(r interface{}) {
	// Same as stdlib http server code. Manually allocate stack trace buffer size
	// to prevent excessively large logs
	const size = 64 << 10
	stacktrace := make([]byte, size)
	stacktrace = stacktrace[:runtime.Stack(stacktrace, false)]
	if _, ok := r.(string); ok {
		log.Printf("Observed a panic: %s\n%s", r, stacktrace)
	} else {
		log.Printf("Observed a panic: %#v (%v)\n%s", r, r, stacktrace)
	}
}

func HandleError(err error, fn func(err error)) {
	if err == nil {
		return
	}
	fn(err)
}

func HandleNil[T any](pointer *T, fn func(notNilPointer *T)) {
	if pointer == nil {
		return
	}
	fn(pointer)
}
