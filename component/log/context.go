package log

import "go.uber.org/zap"

// Context
func Context(items map[string]interface{}) zap.Field {
	return zap.Any("context", items)
}
