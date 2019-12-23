package container

import (
	"go.uber.org/dig"
)

// New 提供自定义容器的接口
func New() *dig.Container {
	return dig.New()
}


