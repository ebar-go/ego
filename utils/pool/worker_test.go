package pool

import (
	"fmt"
	"testing"
)

func BenchmarkWorkerPool_Schedule(b *testing.B) {
	pool := NewWorkerPool(1024)
	for i := 0; i < b.N; i++ {
		pool.Schedule(func() {
			_ = fmt.Sprintf("foo")
		})
	}
}

func BenchmarkGoroutinePool_Schedule(b *testing.B) {
	pool := NewGoroutinePool(1024)
	for i := 0; i < b.N; i++ {
		pool.Schedule(func() {
			_ = fmt.Sprintf("foo")
		})
	}
}
