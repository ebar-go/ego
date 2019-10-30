// response 基于gin的Context,实现响应数据结构体
// 集成全局traceID
package response

import (
	"github.com/ebar-go/ego/library"
	"github.com/gin-gonic/gin"
	"github.com/ebar-go/ego/http/util"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// ErrorItem 错误项
type ErrorItem struct {
	Key   string
	Value string
}

type ResponseInterface interface {
	IsSuccess()
}

// Response 数据结构体
type Response struct {
	context *gin.Context  `json:"-"`
	StatusCode interface{}     `json:"status_code"` // 兼容字符串与int
	Message    string  `json:"message"`
	Data       Data    `json:"data"`
	Meta       Meta    `json:"meta"`
	Errors     []ErrorItem `json:"errors"`
}

// MapResponse 数组类型的数据结构体
type MapResponse struct {
	context *gin.Context  `json:"-"`
	StatusCode interface{}     `json:"status_code"`
	Message    string  `json:"message"`
	Data       []Data    `json:"data"`
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

const (
	StatusOk = 200
)

// 数据对象
type Data map[string]interface{}

// SetContext 设置gin上下文
func (response *Response) SetContext(context *gin.Context)  {
	response.context = context
}

// Json 输出json
func (response *Response) Json()  {
	response.context.JSON(StatusOk, response)
}

// String 输出字符串
func (response *Response) String(format string, values ...interface{})  {
	response.context.String(StatusOk, format, values)
}

// Error 错误输出
func (response *Response) Error(statusCode int, message string)  {
	response.StatusCode = statusCode
	response.Message = message
	response.context.JSON(StatusOk, response)
}

// Success 成功的输出
func (response *Response) Success(data Data)  {
	response.Data = data
	response.context.JSON(StatusOk, response)
}

// IsSuccess 是否已成功
func (response *Response) IsSuccess() bool {
	return response.StatusCode == StatusOk
}

// IsSuccess 是否已成功
func (response *MapResponse) IsSuccess() bool {
	statusCode := ""
	switch reflect.TypeOf(response.StatusCode).Kind() {
	case reflect.Float64:
		statusCode = fmt.Sprintf("%.f", response.StatusCode.(float64))
	case reflect.String:
		statusCode = response.StatusCode.(string)
	}
	return statusCode == strconv.Itoa(StatusOk)
}

// Decode 解析json数据
func Decode(result string) *Response {
	var resp Response
	err := json.Unmarshal([]byte(result), &resp)
	if  err != nil {
		fmt.Println(err)
		return nil
	}

	return &resp
}

// Decode 解析json数据
func DecodeMap(result string) *MapResponse {
	var resp MapResponse
	err := json.Unmarshal([]byte(result), &resp)
	if  err != nil {
		fmt.Println(err)
		return nil
	}

	return &resp
}

