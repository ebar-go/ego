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
