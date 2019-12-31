package json

import (
	"github.com/pquerna/ffjson/ffjson"
)

// Encode json stringify
func Encode(v interface{}) (string, error) {
	buf, err := ffjson.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

// Decode json to interface
func Decode(buf []byte, obj interface{}) error {
	return ffjson.Unmarshal(buf, obj)
}
