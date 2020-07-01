package array

import (
	"fmt"
	"strings"
)

// Array interface
type Array interface {
	Implode(separator string) string
	Length() int
}


// Implode concat items by the given separator
func Implode( items interface{}, separator string) string {
	return strings.Replace(strings.Trim(fmt.Sprint(items), "[]"), " ", separator, -1)
}

// Explode split string with separator
func Explode(str, separator string) StringSlice {
	return String(strings.Split(str, separator))
}
