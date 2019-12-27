package response

import (
	"fmt"
	"github.com/ebar-go/ego/component/pagination"
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/ego/helper"
	"reflect"
	"strconv"
)

// newInstance 实例化response
func newInstance() *Response {
	return &Response{
		StatusCode: 200,
		Message:    "",
		Data:       nil,
		Meta: Meta{
			Trace: Trace{
				RequestId: helper.NewRequestId(),
				TraceId:   trace.GetTraceId(),
			},
		},
		Errors: []ErrorItem{},
	}
}

// Response 数据结构体
type Response struct {
	StatusCode interface{} `json:"status_code"` // 兼容字符串与int
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Meta       Meta        `json:"meta"`
	Errors     []ErrorItem `json:"errors"`
}

// SetStatusCode
func (response *Response) SetStatusCode(code int) {
	response.StatusCode = code
}

func (response *Response) GetMessage() string {
	return response.Message
}

// SetMessage
func (response *Response) SetMessage(message string) {
	response.Message = message
}

// SetErrors set errors
func (response *Response) SetErrors(e []ErrorItem) {
	response.Errors = e
}

// GetData
func (response *Response) GetData() interface{} {
	return response.Data
}

// GetErrors
func (response *Response) GetErrors() []ErrorItem {
	return response.Errors
}

// Trace 跟踪信息
type Trace struct {
	TraceId   string `json:"trace_id"`   // 全局唯一Code
	RequestId string `json:"request_id"` // 当前请求code
}

// Meta 元数据
type Meta struct {
	Trace      Trace                 `json:"trace"`
	Pagination *pagination.Paginator `json:"pagination"` // 分页信息
}

// String 序列化
func (response *Response) String() string {
	resp, _ := helper.JsonEncode(response)
	return resp
}

// SetData 设置数据
func (respone *Response) SetData(data interface{}) {
	respone.Data = data
}

// IsSuccess 是否已成功
func (response *Response) IsSuccess() bool {
	return formatStatusCode(response.StatusCode) == strconv.Itoa(200)
}

// formatStatusCode 格式化
func formatStatusCode(v interface{}) string {
	statusCode := ""
	switch reflect.TypeOf(v).Kind() {
	case reflect.Float64:
		statusCode = fmt.Sprintf("%.f", v.(float64))
	case reflect.String:
		statusCode = v.(string)
	}
	return statusCode
}
