package object

import "io"

// RequestParam 请求参数
type RequestParam struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    io.Reader
}
