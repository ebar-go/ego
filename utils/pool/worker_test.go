package pool

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	w := NewWorker(10, func() {
		log.Println("closed")
	}, func(w *Worker) {
	})

	go func() {
		for i := 0; i < 100000; i++ {
			content := fmt.Sprintf("num: %d", i)
			w.Submit(func() {
				log.Println(content)
			}, false)
			time.Sleep(time.Millisecond * 10)
		}
	}()

	time.Sleep(time.Second * 1)
	w.Stop()
}
