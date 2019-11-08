package helper

import (
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/ebar-go/ego/library"
	"github.com/ebar-go/ego/http/constant"
)


// GetTraceId 获取唯一traceId
func GetTraceId(c *gin.Context) string {
	traceIdInterface, exist := c.Get(constant.TraceID)
	traceId := ""
	if exist == false {
		traceId = c.GetHeader(constant.GatewayTrace)
		if strings.TrimSpace(traceId) == "" {
			traceId = constant.TraceIdPrefix + library.UniqueId()
		}
		c.Set(constant.TraceID, traceId)
	}else {
		traceId = traceIdInterface.(string)
	}

	return traceId
}
