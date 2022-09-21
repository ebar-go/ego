package compress

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	gc = NewGzip()
)

func TestGzipCompressor_Compress(t *testing.T) {
	// test compress
	err := gc.Compress(bytes.NewBuffer([]byte{}), []byte("hello,world"))
	assert.Nil(t, err)

	err = gc.Compress(bytes.NewBuffer([]byte{}), []byte(""))
	assert.Nil(t, err)
}
func TestGzipCompressor_Decompress(t *testing.T) {
	source := []byte("hello,world")
	input := bytes.NewBuffer([]byte{})

	// test compress
	err := gc.Compress(input, source)
	assert.Nil(t, err)

	// test decompress
	output := bytes.NewBuffer([]byte{})
	err = gc.Decompress(output, input.Bytes())
	assert.Nil(t, err)

	// compare decompress result by source
	assert.Equal(t, source, output.Bytes())
}
