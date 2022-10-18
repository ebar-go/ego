package ego

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, New())
}

func TestNewHTTPServer(t *testing.T) {
	assert.NotNil(t, NewHTTPServer(":8080"))
}

func TestNewGRPCServer(t *testing.T) {
	assert.NotNil(t, NewGRPCServer(":8081"))
}
