package middleware

import (
	"testing"
	"fmt"
)

func TestEncodeToken(t *testing.T) {
	tokenStr, err := GetEncodeToken("common-openapi", "WUcLklcyETkhz7ktThMniw6AFseNbrJ6", 600)
	fmt.Println(tokenStr, err)
}
