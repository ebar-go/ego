package compress

import (
	"compress/gzip"
	"sync"
)

// CompressorProvider describes a component that can provider compressors for the std methods.
type CompressorProvider interface {
	// AcquireGzipWriter Returns a *gzip.Writer which needs to be released later.
	// Before using it, call Reset().
	AcquireGzipWriter() *gzip.Writer

	// ReleaseGzipWriter Releases an acquired *gzip.Writer.
	ReleaseGzipWriter(w *gzip.Writer)

	// AcquireGzipReader Returns a *gzip.Reader which needs to be released later.
	AcquireGzipReader() *gzip.Reader

	// ReleaseGzipReader Releases an acquired *gzip.Reader.
	ReleaseGzipReader(r *gzip.Reader)
}

type SyncPoolCompressors struct {
	writerPool *sync.Pool
	readerPool *sync.Pool
}

func NewSyncPoolCompressors() CompressorProvider {
	return &SyncPoolCompressors{
		readerPool: &sync.Pool{
			New: func() any {
				return new(gzip.Reader)
			}},
		writerPool: &sync.Pool{
			New: func() any {
				return new(gzip.Writer)
			}},
	}
}

func (s *SyncPoolCompressors) AcquireGzipWriter() *gzip.Writer {
	return s.writerPool.Get().(*gzip.Writer)
}

func (s *SyncPoolCompressors) ReleaseGzipWriter(w *gzip.Writer) {
	s.writerPool.Put(w)
}

func (s *SyncPoolCompressors) AcquireGzipReader() *gzip.Reader {
	return s.readerPool.Get().(*gzip.Reader)
}

func (s *SyncPoolCompressors) ReleaseGzipReader(r *gzip.Reader) {
	s.readerPool.Put(r)
}
