package pool

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestGoroutine(t *testing.T) {
	pool := NewGoroutinePool()

	for i := 0; i < 1000; i++ {
		pool.Schedule(func() {
			log.Println("test")
		}, true)
	}

	go func() {
		for {
			time.Sleep(time.Second * 1)
			pool.Schedule(func() {
				log.Println("test")
			}, true)
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

func TestFixedGoroutinePool_Schedule(t *testing.T) {

	t.Run("block", func(t *testing.T) {
		pool := NewFixedGoroutinePool()
		for i := 0; i < 1000; i++ {
			content := fmt.Sprintf("content:%d", i)
			pool.Schedule(func() {
				fmt.Println(content)
				time.Sleep(time.Millisecond * 10)
			})
		}

		time.Sleep(time.Second * 3)
	})

	t.Run("nonBlock", func(t *testing.T) {
		pool := NewFixedGoroutinePool(func(options *FixedOptions) {
			options.Block = false
		})
		for i := 0; i < 1000; i++ {
			content := fmt.Sprintf("content:%d", i)
			pool.Schedule(func() {
				fmt.Println(content)
				time.Sleep(time.Millisecond * 10)
			})
		}

		time.Sleep(time.Second * 3)
	})
}
