package mns

import "github.com/ebar-go/ego/helper"

// Params 参数
type Params struct {
	Content interface{} `json:"content"`
	Tag string `json:"tag"`
	TraceId string `json:"trace_id"`
	ReferServiceName string `json:"refer_service_name"`
	Sign string `json:"sign"`
}

// GenerateSign 生成签名
func (params Params) GenerateSign(key string) string {
	return  helper.GetMd5String(params.TraceId + key)
}
