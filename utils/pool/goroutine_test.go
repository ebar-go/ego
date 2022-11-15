package pool

import (
	"log"
	"testing"
)

func TestGoroutine(t *testing.T) {
	pool := NewGoroutinePool()

	for i := 0; i < 1000; i++ {
		pool.Schedule(func() {
			log.Println("test")
		})
	}

	select {}
}
