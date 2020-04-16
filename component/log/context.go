package log

import (
	"github.com/ebar-go/ego/component/trace"
	"go.uber.org/zap"
)

// Context
func Context(items map[string]interface{}) zap.Field {
	items["trace_id"] = trace.GetTraceId()
	return zap.Any("context", items)
}
