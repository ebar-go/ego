package component

type Instance[T any] interface {
	Delegate() T
}
type instance[T any] struct {
	name     string
	delegate T
}

func NewInstance[T any](name string, delegate T) Instance[T] {
	return &instance[T]{name: name, delegate: delegate}
}

func (c *instance[T]) Delegate() T {
	return c.delegate
}
