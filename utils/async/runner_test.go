package async

import (
	"log"
	"testing"
	"time"
)

func TestNewRunner(t *testing.T) {
	runner := NewRunner(func(stop <-chan struct{}) {
		log.Println("first run")
		<-stop
		log.Println("first stop")
	}, func(stop <-chan struct{}) {
		log.Println("other run")
		<-stop
		log.Println("other stop")
	})

	runner.Start()

	time.Sleep(time.Second)
	runner.Stop()
	time.Sleep(time.Second)
}
