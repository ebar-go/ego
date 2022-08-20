package serializer

import (
	"bytes"
	"encoding/json"
	"github.com/ebar-go/ego/rebuild/utils/pool"
)

type Serializer interface {
	Bytes() []byte
	Decode(container interface{}) error
	Release()
}

type Buffer struct {
	*bytes.Buffer
}

func (buffer *Buffer) Decode(container interface{}) error {
	return json.Unmarshal(buffer.Bytes(), container)
}

func (buffer *Buffer) Release() {
	pool.PutByte(buffer.Bytes())
}

func NewBuffer(length int) *Buffer {
	return &Buffer{
		Buffer: bytes.NewBuffer(pool.GetByte(length)),
	}
}
