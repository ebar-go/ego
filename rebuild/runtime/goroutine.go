package runtime

import "log"

func Goroutine(fn func()) {
	GoroutineRecover(fn, func(reason interface{}) {
		log.Println("goroutine recover:", reason)
	})
}

func GoroutineRecover(fn func(), callback func(reason interface{})) {
	defer func() {
		if r := recover(); r != nil {
			callback(r)
		}
	}()
	fn()
}

func HandleCrash() {
	if r := recover(); r != nil {
		log.Println("goroutine crash:", r)
	}
}
