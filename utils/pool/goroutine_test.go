package pool

import (
	"log"
	"testing"
	"time"
)

func TestGoroutine(t *testing.T) {
	pool := NewGoroutinePool()

	for i := 0; i < 1000; i++ {
		pool.Schedule(func() {
			log.Println("test")
		})
	}

	go func() {
		for {
			time.Sleep(time.Second * 1)
			pool.Schedule(func() {
				log.Println("test")
			})
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * 5)
			log.Println("pool state:", len(pool.workers))
		}
	}()

	time.Sleep(time.Minute)
	pool.Stop()
}
