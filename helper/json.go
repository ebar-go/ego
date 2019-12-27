package helper

import json "github.com/pquerna/ffjson/ffjson"

// JsonEncode json序列号
func JsonEncode(v interface{}) (string, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

// JsonDecode json decode
func JsonDecode(buf []byte, obj interface{}) error {
	return json.Unmarshal(buf, obj)
}
