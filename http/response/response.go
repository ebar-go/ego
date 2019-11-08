// response 基于gin的Context,实现响应数据结构体
// 集成全局traceID
package response

import (
	"github.com/ebar-go/ego/library"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"github.com/ebar-go/ego/http/constant"
	"github.com/ebar-go/ego/http/helper"
)

// ErrorItem 错误项
type ErrorItem struct {
	Key   string `json:"key"`
	Value string `json:"error"`
}

// NewErrorItem 实例化错误项
func NewErrorItem(key, msg string) ErrorItem {
	return ErrorItem{Key:key, Value:msg}
}

// IResponse Response接口
type IResponse interface {
	IsSuccess()
}

// 数据对象
type Data map[string]interface{}

// Response 数据结构体
type Response struct {
	StatusCode interface{}     `json:"status_code"` // 兼容字符串与int
	Message    string  `json:"message"`
	Data       interface{}    `json:"data"`
	Meta       Meta    `json:"meta"`
	Errors     []ErrorItem `json:"errors"`
}

// MapResponse 数组类型的数据结构体
type MapResponse struct {
	StatusCode interface{}     `json:"status_code"`
	Message    string  `json:"message"`
	Data       []Data    `json:"data"`
	Meta       Meta    `json:"meta"`
	Errors     []ErrorItem `json:"errors"`
}

// Meta 元数据
type Meta struct {
	TraceId string `json:"trace_id"` // 全局唯一Code
	RequestId string `json:"request_id"` // 当前请求code
	Pagination *library.Pagination `json:"pagination"` // 分页信息
}

// Default 实例化response
func Default() *Response {
	return &Response{
		StatusCode: constant.StatusOk,
		Message: "",
		Data: nil,
		Meta: Meta{
			RequestId: library.UniqueId(),
		},
		Errors: nil,

	}
}

// complete 补全必要参数
func (response *Response) complete()  {
	if &response.Meta == nil {
		response.Meta = Meta{
			RequestId: library.UniqueId(),
		}
	}

}

// Json 输出json
func Json(ctx *gin.Context, response *Response)  {

	response.complete()
	response.Meta.TraceId = helper.GetTraceId(ctx)
	ctx.JSON(constant.StatusOk, response)
}

// Success 成功的输出
func Success(ctx *gin.Context, data interface{})  {
	response := Default()
	response.Data = data
	Json(ctx, response)
}

// Error 错误输出
func Error(ctx *gin.Context, statusCode int, message string)  {
	response := Default()
	response.StatusCode = statusCode
	response.Message = message
	Json(ctx, response)
}


// IsSuccess 是否已成功
func (response *Response) IsSuccess() bool {
	return formatStatusCode(response.StatusCode) == strconv.Itoa(constant.StatusOk)
}

// IsSuccess 是否已成功
func (response *MapResponse) IsSuccess() bool {
	return formatStatusCode(response.StatusCode) == strconv.Itoa(constant.StatusOk)
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

// Decode 解析json数据Response
func Decode(result string) *Response {
	var resp Response
	err := json.Unmarshal([]byte(result), &resp)
	if  err != nil {
		fmt.Println(err)
		return nil
	}

	return &resp
}

// DecodeMap 解析json数据为MapResponse
func DecodeMap(result string) *MapResponse {
	var resp MapResponse
	err := json.Unmarshal([]byte(result), &resp)
	if  err != nil {
		fmt.Println(err)
		return nil
	}

	return &resp
}

