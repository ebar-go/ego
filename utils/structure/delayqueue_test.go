package structure

import (
	"fmt"
	"testing"
	"time"
)

func TestDelayQueue(t *testing.T) {
	dq := NewDelayQueue(10)
	go dq.Poll(func() int64 {
		return time.Now().UnixNano()
	})
	go func() {
		for i := 0; i < 10; i++ {
			dq.Offer(i, time.Now().UnixNano())
			time.Sleep(time.Second * 2)
		}
		dq.Close()
	}()

	for {
		select {
		case item, ok := <-dq.Data:
			if !ok {
				return
			}
			fmt.Println(item)
		default:
			time.Sleep(time.Second)
		}
	}
}
