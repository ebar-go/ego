package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	builder := New()
	builder.Default().SetDefault("foo", "bar")

	res, exist := builder.Default().Get("foo")
	assert.True(t, exist)
	assert.Equal(t, "bar", res)
}
