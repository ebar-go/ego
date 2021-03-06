package curl

import (
	"bytes"
	"encoding/json"
	"io"
	"unsafe"
)

// Response
type Response interface {
	String() string
	Byte() []byte
	BindJson(object interface{}) error
	Reader() io.Reader
}

// response http response wrapper
type response struct {
	body []byte
}

// String return response as string
func (wrap *response) String() string {
	return *(*string)(unsafe.Pointer(&wrap.body))
}

// Byte return response as byte
func (wrap *response) Byte() []byte {
	return wrap.body
}

// BindJson bind json object with pointer
func (wrap *response) BindJson(object interface{}) error {
	return json.Unmarshal(wrap.body, object)
}

// Reader
func (wrap *response) Reader() io.Reader {
	return bytes.NewReader(wrap.body)
}
