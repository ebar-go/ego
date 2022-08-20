package runtime

import (
	"log"
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

func HandleCrash() {
	if r := recover(); r != nil {
		log.Println(r)
	}
}
