// response 基于gin的Context,实现响应数据结构体
// 集成全局traceID
package response

import (
	"github.com/ebar-go/ego/library"
	"github.com/gin-gonic/gin"
	"github.com/ebar-go/ego/http/request"
)

// ErrorItem 错误项
type ErrorItem struct {
	Key   string
	Value string
}

// Response 数据结构体
type Response struct {
	context *gin.Context  `json:"-"`
	StatusCode int     `json:"status_code"`
	Message    string  `json:"message"`
	Data       Data    `json:"data"`
	Meta       Meta    `json:"meta"`
	Errors     []ErrorItem `json:"errors"`
}

type Meta struct {
	TraceId string `json:"trace_id"`
	RequestId string `json:"request_id"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	
}

// Default 实例化response
func Default(context *gin.Context) *Response {
	return &Response{
		context:context,
		StatusCode: 200,
		Message: "",
		Data: nil,
		Meta: Meta{
			TraceId: request.GetTraceId(context),
			RequestId: library.UniqueId(),
		},
		Errors: nil,

	}
}

// 数据对象
type Data map[string]interface{}

// Json 输出json
func (response *Response) Json()  {
	response.context.JSON(200, response)
}

// String 输出字符串
func (response *Response) String(format string, values ...interface{})  {
	response.context.String(200, format, values)
}


