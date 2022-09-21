package compress

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGzip(t *testing.T) {
	compressor := NewGzip()
	assert.NotNil(t, compressor)
}
