// response 基于gin的Context,实现响应数据结构体
// 集成全局traceID
package response

import (
	"github.com/ebar-go/ego/http/pagination"
	"github.com/ebar-go/ego/utils/json"
	"github.com/gin-gonic/gin"
	"reflect"
)

// IResponse Response接口
type IResponse interface {
	// 序列话
	String() string

	// 判断是否成功
	IsSuccess() bool

	// 设置数据项
	SetData(data interface{})

	// 设置状态码
	SetStatusCode(code int)

	// 设置信息
	SetMessage(message string)

	// 设置错误信息
	SetErrors(errors []ErrorItem)

	// 获取信息
	GetMessage() string

	// 获取数据项
	GetData() interface{}

	// 获取错误项
	GetErrors() []ErrorItem
}

// 数据对象
type Data map[string]interface{}

// Default 实例化response
func Default() IResponse {
	return newInstance()
}

// Paginate 分页输出
// formatMap 是否将data项格式化为数组
func Paginate(ctx *gin.Context, data interface{}, paginate *pagination.Paginator, formatMap bool) {
	resp := newInstance()

	v := reflect.ValueOf(data)
	if formatMap && v.IsNil() {
		resp.SetData([]interface{}{})
	} else {
		resp.SetData(data)
	}

	resp.Meta.Pagination = paginate
	Json(ctx, resp)
}

// Json 输出json
func Json(ctx *gin.Context, response IResponse) {
	ctx.JSON(200, response)
}

// Success 成功的输出
func Success(ctx *gin.Context, data interface{}) {
	response := Default()
	response.SetData(data)
	Json(ctx, response)
}

// Error 错误输出
func Error(ctx *gin.Context, statusCode int, message string) {
	response := Default()
	response.SetStatusCode(statusCode)
	response.SetMessage(message)
	Json(ctx, response)
}

// Decode 解析json数据Response
func Decode(result string) IResponse {
	resp := Default()

	if err := json.Decode([]byte(result), resp);err != nil {
		return nil
	}

	return resp
}
