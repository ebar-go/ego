package compressor

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGzip(t *testing.T) {
	compressor := NewGzip()
	assert.NotNil(t, compressor)

	source := []byte("hello,world")
	input := bytes.NewBuffer([]byte{})

	// test compress
	err := compressor.Compress(input, source)
	assert.Nil(t, err)

	// test decompress
	output := bytes.NewBuffer([]byte{})
	err = compressor.Decompress(output, input.Bytes())
	assert.Nil(t, err)

	// compare decompress result by source
	assert.Equal(t, source, output.Bytes())
}
