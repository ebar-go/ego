// response 基于gin的Context,实现响应数据结构体
// 集成全局traceID
package response

import (
	"github.com/ebar-go/ego/library"
	"github.com/gin-gonic/gin"
	"github.com/ebar-go/ego/http/util"
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
	Pagination library.Pagination `json:"pagination"`
}



// Default 实例化response
func Default(context *gin.Context) *Response {
	return &Response{
		context:context,
		StatusCode: 200,
		Message: "",
		Data: nil,
		Meta: Meta{
			TraceId: util.GetTraceId(context),
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

// Error 错误输出
func (response *Response) Error(statusCode int, message string)  {
	response.StatusCode = statusCode
	response.Message = message
	response.context.JSON(200, response)
}

// Success 成功的输出
func (response *Response) Success(data Data)  {
	response.Data = data
	response.context.JSON(200, response)
}


