package curl

import (
	"bytes"
	"github.com/ebar-go/egu"
	"io"
)

// response http response wrapper
type response struct {
	body []byte
}

// String return response as string
func (wrap *response) String() (string) {
	return egu.Byte2Str(wrap.body)
}


// Byte return response as byte
func (wrap *response) Byte() ([]byte) {
	return wrap.body
}

// BindJson bind json object with pointer
func (wrap *response) BindJson(object interface{}) error {
	return egu.JsonDecode(wrap.body, object)
}

// Reader
func (wrap *response) Reader() io.Reader {
	return bytes.NewReader(wrap.body)
}
