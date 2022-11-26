package structure

import "sync"

type Singleton[T any] struct {
	once        sync.Once
	instance    T
	constructor func() T
}

func (s *Singleton[T]) Get() T {
	s.once.Do(func() {
		s.instance = s.constructor()
	})
	return s.instance
}

func NewSingleton[T any](constructor func() T) *Singleton[T] {
	return &Singleton[T]{
		once:        sync.Once{},
		constructor: constructor,
	}
}
