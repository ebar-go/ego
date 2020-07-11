package event

var d = &dispatcher{items: map[string][]Listener{}}

// dispatcher
type dispatcher struct {
	items map[string][]Listener
}

// Register
func Register(eventName string, listener Listener) {
	if Has(eventName) {
		d.items[eventName] = append(d.items[eventName], listener)
	} else {
		d.items[eventName] = []Listener{listener}
	}
}

// Listen
func Listen(eventName string, handler Handler) {
	Register(eventName, Listener{
		Mode:    Sync,
		Handler: handler,
	})
}

// Has
func Has(eventName string) bool {
	_, ok := d.items[eventName]
	return ok
}

// Trigger
func Trigger(eventName string, params interface{}) {
	if !Has(eventName) {
		return
	}

	for _, listener := range d.items[eventName] {
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
