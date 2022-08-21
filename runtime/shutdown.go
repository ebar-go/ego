package runtime

import (
	"os"
	"os/signal"
	"syscall"
)

// shutdown 关闭服务
func Shutdown(callback func()) {
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-c
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if callback != nil {
				callback()
			}

			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

// WaitClose is a helper function that waits for the channel to close and invoke the callback
func WaitClose(stop <-chan struct{}, onClose func()) {
	for {
		select {
		case <-stop:
			onClose()
		default:

		}
	}
}
