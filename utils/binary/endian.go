package binary

import (
	"bytes"
)

type Endian interface {
	Int8(b []byte) int8
	PutInt8(b []byte, v int8)

	Int16(b []byte) int16
	PutInt16(b []byte, v int16)

	Int32(b []byte) int32
	PutInt32(b []byte, v int32)

	String(b []byte) string
	PutString(b []byte, v string)
}

// BigEndian big endian.
func BigEndian() Endian {
	return defaultBigEndian
}

var (
	defaultBigEndian = bigEndian{}
)

type bigEndian struct{}

func (bigEndian) Int8(b []byte) int8 { return int8(b[0]) }

func (bigEndian) PutInt8(b []byte, v int8) {
	b[0] = byte(v)
}

func (bigEndian) Int16(b []byte) int16 { return int16(b[1]) | int16(b[0])<<8 }

func (bigEndian) PutInt16(b []byte, v int16) {
	b[0] = byte(v >> 8)
	b[1] = byte(v)
}

func (bigEndian) Int32(b []byte) int32 {
	return int32(b[3]) | int32(b[2])<<8 | int32(b[1])<<16 | int32(b[0])<<24
}

func (bigEndian) PutInt32(b []byte, v int32) {
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}

func (bigEndian) String(b []byte) string {
	index := bytes.IndexByte(b, 0)
	if index < 0 {
		return ""
	}
	buf := b[0:index]
	return string(buf)
}
func (bigEndian) PutString(b []byte, v string) {
	destLen := len(b)
	source := []byte(v)
	sourceLen := len(source)
	for i := 0; i < sourceLen; i++ {
		if i < destLen {
			b[i] = source[i]
		}
	}
}
