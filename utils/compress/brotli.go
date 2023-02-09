package compress

import (
	"bytes"
	"github.com/andybalholm/brotli"
	"io"
)

type BrotliCompressor struct {
}

func (compressor BrotliCompressor) Compress(dst io.Writer, src []byte) (err error) {
	if len(src) == 0 {
		return
	}

	w := brotli.NewWriter(dst)
	_, err = w.Write(src)

	return
}

func (compressor BrotliCompressor) Decompress(dst io.Writer, src []byte) (err error) {
	r := brotli.NewReader(bytes.NewReader(src))

	_, err = io.Copy(dst, r)

	return
}
