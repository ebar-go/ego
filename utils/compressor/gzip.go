package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"
)

type gzipCompress struct {
	rp *sync.Pool
	wp *sync.Pool
}

func (c *gzipCompress) Compress(dst io.Writer, src []byte) (err error) {
	// not compress empty bytes.
	if len(src) == 0 {
		return
	}

	w := c.wp.Get().(*gzip.Writer)
	w.Reset(dst)
	defer c.wp.Put(w)

	return c.compress(w, src)
}

func (c *gzipCompress) Decompress(dst io.Writer, src []byte) (err error) {
	r := c.rp.Get().(*gzip.Reader)
	defer c.rp.Put(r)

	if err = r.Reset(bytes.NewReader(src)); err != nil {
		return
	}
	_, err = io.Copy(dst, r)
	return
}

func (c *gzipCompress) compress(w *gzip.Writer, src []byte) error {
	if _, err := w.Write(src); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}
