package encrypt

import "encoding/base64"

// Base64Decode decode base64 string
func Base64Decode(encoded string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil
	}
	return decoded
}

// Base64Encode return base64 string
func Base64Encode(source []byte) string {
	return base64.StdEncoding.EncodeToString(source)
}
