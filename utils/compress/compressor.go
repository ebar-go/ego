package compress

import (
	"io"
)

// Compressor represents compress interface
type Compressor interface {
	Compress(dst io.Writer, src []byte) (err error)
	Decompress(dst io.Writer, src []byte) (err error)
}
