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

func BenchmarkCompress(b *testing.B) {
	//BenchmarkCompress/gzip
	//BenchmarkCompress/gzip-8                 2603155               488.5 ns/op
	//             309 B/op          0 allocs/op
	//BenchmarkCompress/brotli
	//BenchmarkCompress/brotli-8                 40953             30429 ns/op
	//             258 B/op          0 allocs/op
	// gzip is much better than brotli in golang
	target := []byte("Don't communicate by sharing memory, share memory by communicating.")
	b.Run("gzip", func(b *testing.B) {
		b.ReportAllocs()
		compressor := NewGzipCompressor()
		input := bytes.NewBuffer([]byte{})
		for i := 0; i < b.N; i++ {
			compressor.Compress(input, target)
		}

	})
	b.Run("brotli", func(b *testing.B) {
		b.ReportAllocs()
		compressor := NewBrotliCompressor()

		input := bytes.NewBuffer([]byte{})
		for i := 0; i < b.N; i++ {
			compressor.Compress(input, target)
		}

	})

}
