package component

// Component represents a component
type Component interface {
	Name() string
}

type Named string

func (n Named) Name() string {
	return string(n)
}
