package timewheel

import (
	"fmt"
	"testing"
	"time"
)

func TestDelayQueue(t *testing.T) {
	queue := NewDelayQueue(10)
	queue.Offer("a", 1)
	queue.Offer("b", 2)
	queue.Offer("c", 3)

	exitC := make(chan struct{})
	go func() {
		queue.Poll(exitC, func() int64 {
			return time.Now().UnixMilli()
		})
	}()

	go func() {
		for {
			select {
			case <-exitC:
				return
			case item := <-queue.Channel():
				fmt.Println(item)

			}
		}
	}()

	time.Sleep(time.Second)
	exitC <- struct{}{}

}
