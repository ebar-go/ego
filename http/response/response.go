package response

import (
	"fmt"
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/ego/http/pagination"
	uuid "github.com/satori/go.uuid"
)

func newRequestId() string {
	return fmt.Sprintf("request:%s", uuid.NewV4().String())
}

// Response 数据结构体
type response struct {
	Code    interface{} `json:"code"` // 兼容字符串与int
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
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

// Reset response
func (r *response) Reset() {
	r.Code = 0
	r.Message = "success"
	r.Data = nil
	r.Meta = Meta{
		Trace: Trace{
			RequestId: newRequestId(),
			TraceId:   trace.Get(),
		},
	}
}

// 数据对象
type Data map[string]interface{}
