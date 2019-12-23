package mns

// Params 参数
type Params struct {
	Content          interface{} `json:"content"`
	Tag              string      `json:"tag"`
	TraceId          string      `json:"trace_id"`
	ReferServiceName string      `json:"refer_service_name"`
	Sign             string      `json:"sign"`
}
