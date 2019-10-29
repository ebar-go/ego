package util

import (
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/ebar-go/ego/library"
)

const(
	TraceID = "trace_id" // 全局trace_id
	GatewayTrace = "gateway-trace" // 网关trace
)

// 获取唯一traceId
func GetTraceId(c *gin.Context) string {
	traceIdInterface, exist := c.Get(TraceID)
	traceId := ""
	if exist == false {
		traceId = c.GetHeader(GatewayTrace)
		if strings.TrimSpace(traceId) == "" {
			traceId = library.UniqueId()
		}
		c.Set(TraceID, traceId)
	}else {
		traceId = traceIdInterface.(string)
	}

	return traceId
}
