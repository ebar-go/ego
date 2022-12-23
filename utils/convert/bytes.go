package convert

import (
	"reflect"
	"unsafe"
)

func String2Byte(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	slice := reflect.SliceHeader{Data: sh.Data, Len: sh.Len, Cap: sh.Len}
	return *(*[]byte)(unsafe.Pointer(&slice))
}

func Byte2String(p []byte) string {
	return *(*string)(unsafe.Pointer(&p))
}
