package curl

import (
	"github.com/ebar-go/ego/utils/conv"
	"github.com/ebar-go/ego/utils/json"
)

// response http response wrapper
type response struct {
	body []byte
}

// String return response as string
func (wrap *response) String() (string) {
	return conv.Byte2Str(wrap.body)
}


// Byte return response as byte
func (wrap *response) Byte() ([]byte) {
	return wrap.body
}

// BindJson bind json object with pointer
func (wrap *response) BindJson(object interface{}) error {
	return json.Decode(wrap.body, object)
}
