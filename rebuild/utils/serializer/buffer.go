package serializer

import (
	"bytes"
	"encoding/json"
	"github.com/ebar-go/ego/rebuild/utils/pool"
)

// Serializer provide bytes serialization to pointer
type Serializer interface {
	// Bytes returns the origin bytes data.
	Bytes() []byte

	// Decode serialize bytes to pointer value.
	Decode(container interface{}) error

	// Release put bytes to pool, avoiding memory leak and improve GC performance
	Release()

	// DecodeAndRelease execute Decode and Release function
	DecodeAndRelease(container interface{}) error
}

// BufferSerializer implements Serializer interface by bytes.Buffer.
type BufferSerializer struct {
	*bytes.Buffer
}

func (buffer *BufferSerializer) Decode(container interface{}) error {
	return json.Unmarshal(buffer.Bytes(), container)
}

func (buffer *BufferSerializer) DecodeAndRelease(container interface{}) error {
	err := buffer.Decode(container)
	buffer.Release()
	return err
}

func (buffer *BufferSerializer) Release() {
	pool.PutByte(buffer.Bytes())
}

func NewBuffer(length int) *BufferSerializer {
	return &BufferSerializer{
		Buffer: bytes.NewBuffer(pool.GetByte(length)),
	}
}
