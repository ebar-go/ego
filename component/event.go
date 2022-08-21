package component

import (
	"github.com/ebar-go/ego/runtime"
	"sync"
)

// Event
type Event struct {
	// name
	Name string
	// event params
	Params interface{}
}

const (
	Sync  = 0
	Async = 1
)

// Listener
type Listener struct {
	Mode    int
	Handler Handler
}

// Handler process event
type Handler func(ev Event)

type EventDispatcher struct {
	Named
	items map[string][]Listener
	rmw   sync.RWMutex
}

// Register
func (instance *EventDispatcher) Register(eventName string, listener Listener) {
	instance.rmw.Lock()
	defer instance.rmw.Unlock()
	listeners, ok := instance.items[eventName]
	if !ok {
		// 预定义数组的长度为10
		listeners = make([]Listener, 10)
	}
	listeners = append(listeners, listener)
	instance.items[eventName] = listeners
}

// Listen register a sync event
func (instance *EventDispatcher) Listen(eventName string, handler Handler) {
	instance.Register(eventName, Listener{
		Mode:    Sync,
		Handler: handler,
	})
}

// Has return event exist
func (instance *EventDispatcher) Has(eventName string) bool {
	instance.rmw.RLock()
	defer instance.rmw.RUnlock()
	_, ok := instance.items[eventName]
	return ok
}

// Trigger make event trigger with given name and params
func (instance *EventDispatcher) Trigger(eventName string, params interface{}) {
	instance.rmw.RLock()
	defer instance.rmw.RUnlock()
	listeners, ok := instance.items[eventName]
	if !ok {
		return
	}

	for _, listener := range listeners {
		if listener.Mode == Sync {
			listener.Handler(Event{
				Name:   eventName,
				Params: params,
			})
		} else {
			runtime.Goroutine(func() {
				listener.Handler(Event{
					Name:   eventName,
					Params: params,
				})
			})
		}

	}
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{Named: componentEventDispatcher, items: make(map[string][]Listener, 16)}
}
