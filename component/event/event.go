package event

import (
	"github.com/ebar-go/ego/utils/structure"
	"sync"
)

// Handler process event
type Handler func(param any)

type Instance struct {
	items *structure.ConcurrentMap[string, []Handler]
	rmw   sync.RWMutex
}

// Listen register a sync event
func (instance *Instance) Listen(eventName string, handler Handler) {
	handlers, ok := instance.items.Get(eventName)
	if !ok {
		// 预定义数组的cap为10
		handlers = make([]Handler, 0, 10)
	}
	handlers = append(handlers, handler)
	instance.items.Set(eventName, handlers)
}

// Has return event exist
func (instance *Instance) Has(eventName string) bool {
	_, ok := instance.items.Get(eventName)
	return ok
}

// Trigger make event trigger with given name and params
func (instance *Instance) Trigger(eventName string, param any) {
	handlers, ok := instance.items.Get(eventName)
	if !ok {
		return
	}

	for _, handler := range handlers {
		handler(param)
	}
}

func (instance *Instance) TriggerAsync(eventName string, param any) {
	handlers, ok := instance.items.Get(eventName)
	if !ok {
		return
	}

	for _, handler := range handlers {
		go handler(param)
	}
}

func New() *Instance {
	return &Instance{items: structure.NewConcurrentMap[string, []Handler]()}
}
