package compress

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompress(t *testing.T) {
	source := []byte("hello,world")

	t.Run("gzip", func(t *testing.T) {
		compressor := NewGzipCompressor()
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
	})

	t.Run("brotli", func(t *testing.T) {
		compressor := NewBrotliCompressor()
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
	})

}