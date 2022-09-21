package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/ebar-go/ego/runtime"
	"io"
	"sync"
)

type GzipCompressor struct {
	rp *sync.Pool
	wp *sync.Pool
}

func (c *GzipCompressor) Compress(dst io.Writer, src []byte) (err error) {
	// not compress empty bytes.
	if len(src) == 0 {
		return
	}

	w := c.wp.Get().(*gzip.Writer)
	w.Reset(dst)
	defer c.wp.Put(w)

	return runtime.Call(func() error {
		_, err := w.Write(src)
		return err
	}, w.Flush, w.Close)
}

func (c *GzipCompressor) Decompress(dst io.Writer, src []byte) (err error) {
	r := c.rp.Get().(*gzip.Reader)
	defer c.rp.Put(r)

	if err = r.Reset(bytes.NewReader(src)); err != nil {
		return
	}
	_, err = io.Copy(dst, r)
	return
}
