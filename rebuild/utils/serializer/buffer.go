package serializer

import (
	"bytes"
	"encoding/json"
	"github.com/ebar-go/ego/rebuild/utils/pool"
)

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
	buffer := bytes.NewBuffer(pool.GetByte(length))
	buffer.Reset()
	return &BufferSerializer{
		Buffer: buffer,
	}
}
