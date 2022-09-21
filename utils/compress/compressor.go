package compress

import (
	"io"
)

type Compressor interface {
	Compress(dst io.Writer, src []byte) (err error)
	Decompress(dst io.Writer, src []byte) (err error)
}

func NewGzip() Compressor {
	return &GzipCompressor{provider: CurrentCompressorProvider()}
}
