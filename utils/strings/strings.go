package strings

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
	"strings"
)

import (
	"fmt"
)

//Md5 return the encrypt string by md5 algorithm
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Hash return hash string
func Hash(s string) string {
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(s))
	result := Sha1Inst.Sum([]byte(""))
	return fmt.Sprintf("%x", result)
}

// UUID return unique id
func UUID() string {
	return uuid.NewV4().String()
}

// DecodeBase64
func DecodeBase64(encoded string) string {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return ""
	}
	return string(decoded)
}

// EncodeBase64 return base64 string
func EncodeBase64(source []byte) string {
	return base64.StdEncoding.EncodeToString(source)
}

// Implode concat items by the given separator
func Implode(separator string, items interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(items), "[]"), " ", separator, -1)
}

// Explode split string with separator
func Explode(str, separator string) []string {
	return strings.Split(str, separator)
}

// Default return defaultV if v is empty
func Default(v, defaultV string) string {
	if v == "" {
		return defaultV
	}

	return v
}

func ToBool(b string) bool {
	if b == "1" || "true" == strings.ToLower(b) {
		return true
	}

	return false
}