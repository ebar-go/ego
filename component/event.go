package component

import (
	"fmt"
	"sync"
)

// Handler process event
type Handler func(param any)
type EventDispatcher struct {
	Named
	items map[string][]Handler
	rmw   sync.RWMutex
}

// Listen register a sync event
func (instance *EventDispatcher) Listen(eventName string, handler Handler) {
	instance.rmw.Lock()
	defer instance.rmw.Unlock()
	handlers, ok := instance.items[eventName]
	if !ok {
		// 预定义数组的长度为10
		handlers = make([]Handler, 0, 10)
	}
	handlers = append(handlers, handler)
	instance.items[eventName] = handlers
}

// Has return event exist
func (instance *EventDispatcher) Has(eventName string) bool {
	instance.rmw.RLock()
	defer instance.rmw.RUnlock()
	_, ok := instance.items[eventName]
	return ok
}

// Trigger make event trigger with given name and params
func (instance *EventDispatcher) Trigger(eventName string, param any) {
	fmt.Printf("dispatcher: %p, event: %s\n", instance, eventName)
	instance.rmw.RLock()
	defer instance.rmw.RUnlock()
	handlers, ok := instance.items[eventName]
	if !ok {
		return
	}

	for _, handler := range handlers {
		handler(param)
	}
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{Named: componentEventDispatcher, items: make(map[string][]Handler)}
}

// ListenEvent listen with type parameters
func ListenEvent[T any](eventName string, handler func(param T)) {
	dispatcher := Provider().EventDispatcher()
	dispatcher.Listen(eventName, func(param any) {
		data, ok := param.(T)
		if !ok {
			return
		}
		handler(data)
	})
	fmt.Printf("dispatcher: %p, event: %s\n", dispatcher, eventName)
}

type Event[T any] struct {
	Name string
}

// NewEvent creates a new Event with the given name.
func NewEvent[T any](name string) Event[T] {
	return Event[T]{Name: name}
}

// Bind binds handler
func (e Event[T]) Bind(handler func(param T)) {
	Provider().EventDispatcher().Listen(e.Name, func(param any) {
		data, ok := param.(T)
		if !ok {
			return
		}
		handler(data)
	})
}
