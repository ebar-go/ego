package pool

import (
	"fmt"
	"testing"
	"time"
)

func TestGoroutine(t *testing.T) {
	t.Run("block", func(t *testing.T) {
		pool := NewGoroutinePool()
		for i := 0; i < 1000; i++ {
			content := fmt.Sprintf("content:%d", i)
			pool.Schedule(func() {
				fmt.Println(content)
				time.Sleep(time.Millisecond * 10)
			})
		}

		time.Sleep(time.Second * 3)
		pool.Stop()
	})

	t.Run("nonBlock", func(t *testing.T) {
		pool := NewGoroutinePool(func(options *Options) {
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
		pool.Stop()
	})
}

func TestFixedGoroutinePool(t *testing.T) {

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
		pool.Stop()
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
		pool.Stop()
	})
}
