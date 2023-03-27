package compress

import (
	"bytes"
	"github.com/andybalholm/brotli"
	"github.com/ebar-go/ego/utils/pool"
	"github.com/ebar-go/ego/utils/runtime"
	"io"
)

type BrotliCompressor struct {
	wp pool.Provider[*brotli.Writer]
	rp pool.Provider[*brotli.Reader]
}

func NewBrotliCompressor() *BrotliCompressor {
	return &BrotliCompressor{
		wp: pool.NewSyncPoolProvider[*brotli.Writer](func() interface{} {
			return brotli.NewWriter(nil)
		}),
		rp: pool.NewSyncPoolProvider[*brotli.Reader](func() interface{} {
			return brotli.NewReader(nil)
		}),
	}
}

func (compressor BrotliCompressor) Compress(dst io.Writer, src []byte) (err error) {
	if len(src) == 0 {
		return
	}

	w := compressor.wp.Acquire()
	defer compressor.wp.Release(w)
	w.Reset(dst)
	return runtime.Call(func() error {
		_, err := w.Write(src)
		return err
	}, w.Flush)
}

func (compressor BrotliCompressor) Decompress(dst io.Writer, src []byte) (err error) {
	r := compressor.rp.Acquire()
	defer compressor.rp.Release(r)
	err = r.Reset(bytes.NewReader(src))
	if err != nil {
		return
	}

	_, err = io.Copy(dst, r)
	return
}
