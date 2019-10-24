package library

import "encoding/json"

// JsonEncode json序列号
func JsonEncode(v interface{}) (string, error) {
	bytes , err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
