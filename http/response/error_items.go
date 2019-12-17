package response

// IErrorItems 错误项接口
type IErrorItems interface {
	// 添加错误项
	Push(key, msg string)

	// 是否为空
	IsEmpty() bool
}

// ErrorItem 错误项
type ErrorItem struct {
	Key   string `json:"key"`
	Value string `json:"error"`
}

// ErrorItems 错误项
type ErrorItems struct {
	items []ErrorItem
}

// Push 添加错误项
func (e *ErrorItems) Push(key, msg string) {
	e.items = append(e.items, ErrorItem{Key: key, Value: msg})
}

// IsEmpty 查看错误项是否为空
func (e *ErrorItems) IsEmpty() bool {
	return len(e.items) == 0
}

