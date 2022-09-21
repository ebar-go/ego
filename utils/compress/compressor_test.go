package compress

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefault(t *testing.T) {
	c := Default()
	assert.Equal(t, GzipInstance(), c)
}

func TestCompress(t *testing.T) {
	source := []byte("hello,world")
	input := bytes.NewBuffer([]byte{})

	// test compress
	err := Compress(input, source)
	assert.Nil(t, err)

	// test decompress
	output := bytes.NewBuffer([]byte{})
	err = Decompress(output, input.Bytes())
	assert.Nil(t, err)

	// compare decompress result by source
	assert.Equal(t, source, output.Bytes())
}
