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

func TestDefault(t *testing.T) {
	cache := Default()
	cache.SetDefault("foo", "bar")

	res, exist := cache.Get("foo")
	assert.True(t, exist)
	assert.Equal(t, "bar", res)
}

func TestBuild(t *testing.T) {
	cache := Build()
	cache.SetDefault("foo", "bar")
	res, exist := cache.Get("foo")
	assert.True(t, exist)
	assert.Equal(t, "bar", res)
}
