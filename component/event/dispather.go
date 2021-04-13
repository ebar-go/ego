package event

import "sync"

// dispatcher
type dispatcher struct {
	items map[string][]Listener
	rmw   sync.RWMutex
}

var (
	// 初始化事件分发器，提前给map分配空间，减少因数组扩容带来的消耗
	instance = &dispatcher{items: make(map[string][]Listener, 100), rmw: sync.RWMutex{}}
)

// Register
func Register(eventName string, listener Listener) {
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
func Listen(eventName string, handler Handler) {
	Register(eventName, Listener{
		Mode:    Sync,
		Handler: handler,
	})
}

// Has return event exist
func Has(eventName string) bool {
	instance.rmw.RLock()
	defer instance.rmw.RUnlock()
	_, ok := instance.items[eventName]
	return ok
}

// Trigger
func Trigger(eventName string, params interface{}) {
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
			go listener.Handler(Event{
				Name:   eventName,
				Params: params,
			})
		}

	}
}
