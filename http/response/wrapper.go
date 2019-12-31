// response 基于gin的Context,实现响应数据结构体
// 集成全局traceID
package response

import (
	"github.com/ebar-go/ego/http/pagination"
	"github.com/ebar-go/ego/utils/json"
	"github.com/gin-gonic/gin"
	"reflect"
)

// Wrapper include context
type Wrapper struct {
	ctx *gin.Context
}

// WrapperContext
func WrapperContext(ctx *gin.Context) *Wrapper {
	return &Wrapper{ctx:ctx}
}

// 数据对象
type Data map[string]interface{}

// Paginate 分页输出
// formatMap 是否将data项格式化为数组
func (wrapper *Wrapper) Paginate( data interface{}, paginate *pagination.Paginator, formatMap bool) {
	resp := New()

	v := reflect.ValueOf(data)
	if formatMap && v.IsNil() {
		resp.Data = []interface{}{}
	} else {
		resp.Data = data
	}

	resp.Meta.Pagination = paginate
	wrapper.Json(resp)
}

// Json 输出json
func (wrapper *Wrapper) Json(response *Response) {
	wrapper.ctx.JSON(200, response)
}

// Success 成功的输出
func (wrapper *Wrapper) Success( data interface{}) {
	response := New()
	response.Data = data
	wrapper.Json(response)
}

// Error 错误输出
func (wrapper *Wrapper) Error( statusCode int, message string) {
	response := New()
	response.StatusCode = statusCode
	response.Message = message
	wrapper.Json(response)
}

// Decode 解析json数据Response
func Decode(result string) *Response {
	resp := New()

	if err := json.Decode([]byte(result), resp);err != nil {
		return nil
	}

	return resp
}
