package compress

import (
	"io"
)

// Compressor represents compress interface
type Compressor interface {
	Compress(dst io.Writer, src []byte) (err error)
	Decompress(dst io.Writer, src []byte) (err error)
}

var (
	// Default alias GzipInstance
	Default    = GzipInstance
	Compress   = Default().Compress
	Decompress = Default().Decompress
)
