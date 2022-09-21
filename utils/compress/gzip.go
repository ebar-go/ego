package compress

import (
	"bytes"
	"github.com/ebar-go/ego/runtime"
	"io"
	"sync"
)

type GzipCompressor struct {
	provider CompressorProvider
}

func (c *GzipCompressor) Compress(dst io.Writer, src []byte) (err error) {
	// not compress empty bytes.
	if len(src) == 0 {
		return
	}

	w := c.provider.AcquireGzipWriter()
	w.Reset(dst)
	defer c.provider.ReleaseGzipWriter(w)

	return runtime.Call(func() error {
		_, err := w.Write(src)
		return err
	}, w.Flush, w.Close)
}

func (c *GzipCompressor) Decompress(dst io.Writer, src []byte) (err error) {
	r := c.provider.AcquireGzipReader()
	defer c.provider.ReleaseGzipReader(r)

	if err = r.Reset(bytes.NewReader(src)); err != nil {
		return
	}
	_, err = io.Copy(dst, r)
	return
}

var gzipCompressorInstance struct {
	once     sync.Once
	instance *GzipCompressor
}

func GzipInstance() *GzipCompressor {
	gzipCompressorInstance.once.Do(func() {
		gzipCompressorInstance.instance = &GzipCompressor{provider: CurrentCompressorProvider()}
	})
	return gzipCompressorInstance.instance
}
