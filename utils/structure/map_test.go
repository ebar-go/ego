package structure

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConcurrentMap(t *testing.T) {
	m := NewConcurrentMap[string, string]()

	t.Run("testLen", func(t *testing.T) {
		m.Set("foo", "bar")
		assert.Equal(t, 1, m.Len())

		m.Del("foo")
		assert.Equal(t, 0, m.Len())
	})
}

func BenchmarkConcurrentMap(b *testing.B) {
	b.Run("mutex", func(pb *testing.B) {
		pb.ReportAllocs()
		m := NewConcurrentMap[string, string]()
		for i := 0; i < pb.N; i++ {
			m.Set("foo", "bar")
		}
	})
	b.Run("lockFree", func(pb *testing.B) {
		pb.ReportAllocs()
		m := NewLockFreeMap[string, string]()
		for i := 0; i < pb.N; i++ {
			m.Set("foo", "bar")
		}
	})

}
