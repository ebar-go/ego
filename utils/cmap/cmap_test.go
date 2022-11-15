package cmap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainer_Del(t *testing.T) {
	container := NewContainer[int, string]()
	key := 1
	container.Set(key, "foo")

	item, exist := container.Get(key)
	assert.Equal(t, "foo", item)
	assert.True(t, exist)

	container.Del(key)
	item, exist = container.Get(key)
	assert.Empty(t, item)
	assert.False(t, exist)

}

func TestContainer_Get(t *testing.T) {
	container := NewContainer[int, string]()
	item, exist := container.Get(1)
	assert.Empty(t, item)
	assert.False(t, exist)

	container.Set(2, "foo")

	item, exist = container.Get(2)
	assert.Equal(t, "foo", item)
	assert.True(t, exist)

}

func TestContainer_Set(t *testing.T) {
	container := NewContainer[int, string]()
	container.Set(1, "foo")
	container.Set(2, "bar")
}

func TestNewContainer(t *testing.T) {
	container := NewContainer[int, string]()
	assert.NotNil(t, container)
}
