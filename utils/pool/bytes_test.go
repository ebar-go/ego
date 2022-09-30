package pool

import (
	"fmt"
	"sync"
	"testing"
)

func BenchmarkByteSlice(b *testing.B) {
	b.Run("Run.N", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			bs := GetByte(1024)
			PutByte(bs)
		}
	})
	b.Run("Run.Parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bs := GetByte(1024)
				PutByte(bs)
			}
		})
	})

}

func BenchmarkSyncPool(b *testing.B) {
	pool := sync.Pool{New: func() interface{} {
		return make([]byte, 1024)
	}}
	b.Run("Run.N", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			bs := pool.Get()
			pool.Put(bs)
		}
	})
}

func TestByteSlice(t *testing.T) {
	s := make([]byte, 1024)
	copy(s, "hello")
	a := s[:0]
	fmt.Printf("s = %p, %v\n", s, s)
	fmt.Printf("a = %p\n", a)

	s = s[2:4]
	fmt.Printf("s = %p, %v\n", s, s)
}
