package request

import "io"

// Param 请求参数
type Param struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    io.Reader
}

// AddHeader 添加header头
func (param *Param) AddHeader(key, value string) {
	if param.Headers == nil {
		param.Headers = make(map[string]string)
	}

	param.Headers[key] = value
}
