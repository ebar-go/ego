package response

import (
	"fmt"
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/ego/http/pagination"
	"github.com/ebar-go/ego/utils/json"
	"github.com/ebar-go/ego/utils/strings"
	"reflect"
	"strconv"
)

type IResponse interface {

}

// New return response instance
func New() *Response {
	return &Response{
		StatusCode: 200,
		Message:    "",
		Data:       nil,
		Meta: Meta{
			Trace: Trace{
				RequestId: "RequestId" + strings.UUID(),
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

// String stringify response
func (response *Response) String() string {
	resp, _ := json.Encode(response)
	return resp
}

// IsSuccess get response status
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
