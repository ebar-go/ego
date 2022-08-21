package serializer

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

type SerializerBuilder func() Serializer
