package pool

import (
	"math"
	"math/bits"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

var builtinPool BytePool

var (
	// GetByte retrieves a byte slice but uses the global pool
	GetByte = builtinPool.Get
	// PutByte returns a byte slice but uses the global pool
	PutByte = builtinPool.Put
)

// Pool consists of 32 sync.Pool, representing byte slices of length from 0 to 32 in powers of 2.
type BytePool struct {
	pools [32]sync.Pool
}

// Get retrieves a byte slice of the length requested by the caller from pool or allocates a new one.
func (p *BytePool) Get(size int) (buf []byte) {
	if size <= 0 {
		return nil
	}
	if size > math.MaxInt32 {
		return make([]byte, size)
	}
	idx := index(uint32(size))
	ptr, _ := p.pools[idx].Get().(unsafe.Pointer)
	if ptr == nil {
		return make([]byte, 1<<idx)[:size]
	}
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	sh.Data = uintptr(ptr)
	sh.Len = size
	sh.Cap = 1 << idx
	runtime.KeepAlive(ptr)
	return
}

// Put returns the byte slice to the pool.
func (p *BytePool) Put(buf []byte) {
	size := cap(buf)
	if size == 0 || size > math.MaxInt32 {
		return
	}
	idx := index(uint32(size))
	if size != 1<<idx { // this byte slice is not from Pool.Get(), put it into the previous interval of idx
		idx--
	}
	// array pointer
	p.pools[idx].Put(unsafe.Pointer(&buf[:1][0]))
}

func index(n uint32) uint32 {
	return uint32(bits.Len32(n - 1))
}
