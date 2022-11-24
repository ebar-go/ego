package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSchema(t *testing.T) {
	protocol, bind := "test", "127.0.0.1:8080"
	schema := NewSchema(protocol, bind)

	assert.Equal(t, protocol, schema.Protocol)
	assert.Equal(t, bind, schema.Bind)
}

func TestNewHttpSchema(t *testing.T) {
	bind := "127.0.0.1:8080"
	schema := NewHttpSchema(bind)
	assert.Equal(t, HTTP, schema.Protocol)
	assert.Equal(t, bind, schema.Bind)
}

func TestSchema_Address(t *testing.T) {
	schema := NewHttpSchema("127.0.0.1:8080")
	assert.Equal(t, "http://127.0.0.1:8080", schema.Address())
}

func TestSchema_HostPort(t *testing.T) {
	schema := NewHttpSchema("127.0.0.1:8080")
	host, port, err := schema.HostPort()
	assert.Nil(t, err)
	assert.Equal(t, host, "127.0.0.1")
	assert.Equal(t, port, "8080")
}
