package event

import (
	"sync"
)

// Handler process event
type Handler func(param any)
type Dispatcher struct {
	items map[string][]Handler
	rmw   sync.RWMutex
}

// Listen register a sync event
func (instance *Dispatcher) Listen(eventName string, handler Handler) {
	instance.rmw.Lock()
	defer instance.rmw.Unlock()
	handlers, ok := instance.items[eventName]
	if !ok {
		// 预定义数组的cap为10
		handlers = make([]Handler, 0, 10)
	}
	handlers = append(handlers, handler)
	instance.items[eventName] = handlers
}

// Has return event exist
func (instance *Dispatcher) Has(eventName string) bool {
	instance.rmw.RLock()
	defer instance.rmw.RUnlock()
	_, ok := instance.items[eventName]
	return ok
}

// Trigger make event trigger with given name and params
func (instance *Dispatcher) Trigger(eventName string, param any) {
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

func (instance *Dispatcher) TriggerAsync(eventName string, param any) {
	instance.rmw.RLock()
	defer instance.rmw.RUnlock()
	handlers, ok := instance.items[eventName]
	if !ok {
		return
	}

	for _, handler := range handlers {
		go handler(param)
	}
}

func New() *Dispatcher {
	return &Dispatcher{items: make(map[string][]Handler)}
}
