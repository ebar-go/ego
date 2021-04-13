package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_LoadFile(t *testing.T) {
	conf := New()
	err := conf.LoadFile("./example.yaml", "./other.yaml")
	assert.Nil(t, err)
	assert.Equal(t, "ego", conf.Get("other.name"))
	assert.Equal(t, "test", conf.Get("example.name"))
}
