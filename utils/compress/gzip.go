package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/ebar-go/ego/utils/pool"
	"github.com/ebar-go/ego/utils/runtime"
	"io"
)

type GzipCompressor struct {
	wp pool.Provider[*gzip.Writer]
	rp pool.Provider[*gzip.Reader]
}

func (c *GzipCompressor) Compress(dst io.Writer, src []byte) (err error) {
	// not compress empty bytes.
	if len(src) == 0 {
		return
	}

	w := c.wp.Acquire()
	w.Reset(dst)
	defer c.wp.Release(w)

	return runtime.Call(func() error {
		_, err := w.Write(src)
		return err
	}, w.Flush, w.Close)
}

func (c *GzipCompressor) Decompress(dst io.Writer, src []byte) (err error) {
	r := c.rp.Acquire()
	defer c.rp.Release(r)

	if err = r.Reset(bytes.NewReader(src)); err != nil {
		return
	}
	_, err = io.Copy(dst, r)
	return
}

func NewGzipCompressor() *GzipCompressor {
	return &GzipCompressor{
		wp: pool.NewSyncPoolProvider[*gzip.Writer](func() interface{} {
			return &gzip.Writer{}
		}),
		rp: pool.NewSyncPoolProvider[*gzip.Reader](func() interface{} {
			return &gzip.Reader{}
		}),
	}
}
