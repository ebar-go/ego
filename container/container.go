package container

import (
	"go.uber.org/dig"
)

var (
	App *dig.Container
)

func New() *dig.Container {
	return dig.New()
}


