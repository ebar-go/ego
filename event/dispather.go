package event

// Dispatcher
type Dispatcher interface {
	//add event listener
	AddListener(eventType string, listener *Listener)
	//remove event listener
	RemoveListener(eventType string, listener *Listener) bool
	//check listener is or not exist
	HasListener(eventType string) bool
	//dispatch event
	DispatchEvent(event Event) bool

	// trigger event, simple than DispatchEvent
	Trigger(eventType string, params interface{}) bool
}

// eventDispatcher implement of Dispatcher
type eventDispatcher struct {
	savers []*Saver
}

// NewDispatcher new event dispatcher
func NewDispatcher() Dispatcher {
	return new(eventDispatcher)
}

// AddListener add event listener
func (dispatcher *eventDispatcher) AddListener(eventType string, listener *Listener) {
	for _, saver := range dispatcher.savers {
		if saver.Type == eventType {
			saver.Listeners = append(saver.Listeners, listener)
			return
		}
	}

	saver := &Saver{Type: eventType, Listeners: []*Listener{listener}}
	dispatcher.savers = append(dispatcher.savers, saver)
}

// RemoveListener remove listener from dispatcher
func (dispatcher *eventDispatcher) RemoveListener(eventType string, listener *Listener) bool {
	for _, saver := range dispatcher.savers {
		if saver.Type == eventType {
			for i, l := range saver.Listeners {
				if listener == l {
					saver.Listeners = append(saver.Listeners[:i], saver.Listeners[i+1:]...)
					return true
				}
			}
		}
	}
	return false
}

// HasListener check dispatcher is or not include the eventType
func (dispatcher *eventDispatcher) HasListener(eventType string) bool {
	for _, saver := range dispatcher.savers {
		if saver.Type == eventType {
			return true
		}
	}
	return false
}

// DispatchEvent dispatch the given event
func (dispatcher *eventDispatcher) DispatchEvent(event Event) bool {
	for _, saver := range dispatcher.savers {
		if saver.Type == event.Type {
			for _, listener := range saver.Listeners {
				event.Dispatcher = dispatcher
				listener.Handler(event)
			}
			return true
		}
	}
	return false
}

// DispatchEvent dispatch the given event
func (dispatcher *eventDispatcher) Trigger(eventType string, params interface{}) bool {
	event := New(eventType, params)
	return dispatcher.DispatchEvent(event)
}
