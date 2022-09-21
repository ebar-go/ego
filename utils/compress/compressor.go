package compress

import (
	"compress/gzip"
	"io"
	"sync"
)

type Compressor interface {
	Compress(dst io.Writer, src []byte) (err error)
	Decompress(dst io.Writer, src []byte) (err error)
}

func NewGzip() Compressor {
	return &GzipCompressor{
		rp: &sync.Pool{
			New: func() any {
				return new(gzip.Reader)
			}},
		wp: &sync.Pool{New: func() any {
			return new(gzip.Writer)
		}},
	}
}
