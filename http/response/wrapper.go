// r 基于gin的Context,实现响应数据结构体
// 集成全局traceID
package response

import (
	"github.com/ebar-go/ego/http/pagination"
	"github.com/gin-gonic/gin"
	"reflect"
)

// Wrapper include context
type wrapper struct {
	ctx *gin.Context
}

type Abort struct {

}

// WrapContext
func WrapContext(ctx *gin.Context) *wrapper {
	return &wrapper{ctx: ctx}
}

// output output r
func (w *wrapper) output(r *response) {
	w.ctx.JSON(200, r)
	w.abort()
}

// Success 输出成功响应
func (w *wrapper) Success(data interface{}) {
	r := rp.Get()
	r.Data = data

	w.output(r)
}

// abort 中断
func (w *wrapper) abort() {
	panic(&Abort{})
}

// Error 输出错误响应
func (w *wrapper) Error(code int, message string) {
	r := rp.Get()
	r.Code = code
	r.Message = message

	w.output(r)
}

// Paginate 输出分页响应内容
func (w *wrapper) Paginate(data interface{}, pagination *pagination.Paginator) {
	r := rp.Get()
	// 如果data为nil,则默认设置为[]
	v := reflect.ValueOf(data)
	if v.IsNil() {
		data = []interface{}{}
	}
	r.Data = data
	r.Meta.Pagination = pagination

	w.output(r)
}
