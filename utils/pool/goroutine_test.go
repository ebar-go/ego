package pool

import (
	"log"
	"testing"
)

func TestGoroutine(t *testing.T) {
	pool := NewGoroutinePool(100)

	for i := 0; i < 1000; i++ {
		pool.Schedule(func() {
			log.Println("test")
		})
	}

	select {}
}
