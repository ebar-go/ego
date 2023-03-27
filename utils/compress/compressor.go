package compress

import (
	"github.com/ebar-go/ego/utils/structure"
	"io"
)

// Compressor represents compress interface
type Compressor interface {
	Compress(dst io.Writer, src []byte) (err error)
	Decompress(dst io.Writer, src []byte) (err error)
}

var (
	// Default alias GzipInstance
	defaultInstance = structure.NewSingleton(NewGzipCompressor).Get()
	Compress        = defaultInstance.Compress
	Decompress      = defaultInstance.Decompress
)
