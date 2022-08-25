package component

import "go.uber.org/dig"

// DI implements dependencies injection component by *dig.Container
type DI struct {
	Named
	container *dig.Container
}

func (d *DI) Container() *dig.Container {
	return d.container
}

func (d *DI) Inject(constructor interface{}) error {
	return d.container.Provide(constructor)
}

func (d *DI) Invoke(constructor interface{}) error {
	return d.container.Invoke(constructor)
}

func NewDI() *DI {
	return &DI{Named: componentDI, container: dig.New()}
}
