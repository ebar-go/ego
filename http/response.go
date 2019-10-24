// response 基于gin的Context,实现响应数据结构体
// 集成全局traceID
package http

import (
	"encoding/json"
	"github.com/ebar-go/ego/library"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"github.com/pkg/errors"
	"net/http"
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

// DefaultResponse 实例化response
func DefaultResponse(context *gin.Context) *Response {
	return &Response{
		context:context,
		StatusCode: 200,
		Message: "",
		Data: nil,
		Meta: Meta{
			TraceId: GetTraceId(context),
			RequestId: library.UniqueId(),
		},
		Errors: nil,

	}
}

// 数据对象
type Data map[string]interface{}

// 将响应数据体字符串化
func (object Response) stringify() string {
	str, _ := json.Marshal(object)
	return string(str)
}

// Json 输出json
func (response *Response) Json()  {
	response.context.JSON(200, response)
}

// String 输出字符串
func (response *Response) String(format string, values ...interface{})  {
	response.context.String(200, format, values)
}

// 将response序列化
func  StringifyResponse(response *http.Response) (string, error) {
	if response == nil {
		return "", errors.New("没有响应数据")
	}

	if response.StatusCode != 200 {
		return "", errors.New("非200的上游返回")
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.WithMessage(err, "读取响应数据失败:")
	}

	// 关闭响应
	defer func() {
		response.Body.Close()
	}()

	return string(data), nil
}

