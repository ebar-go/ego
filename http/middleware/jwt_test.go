package middleware

import (
	"testing"
	"fmt"
)

func TestEncodeToken(t *testing.T) {
	tokenStr, err := GetEncodeToken("common-openapi", "aa", 600)
	fmt.Println(tokenStr, err)
}
