package event

import "fmt"

// Event
type Event struct {
	// dispatcher
	Dispatcher Dispatcher
	// name
	Type string
	// event params
	Params interface{}
}

// Saver
type Saver struct {
	// event name
	Type      string
	// event listeners
	Listeners []*Listener
}

// NewEvent
func New(eventType string, params interface{}) Event {
	e := Event{Type: eventType, Params: params}
	return e
}

// Clone Event
func (event *Event) Clone() *Event {
	e := new(Event)
	e.Type = event.Type
	e.Dispatcher = e.Dispatcher
	return e
}

// ToString stringify event
func (event *Event) ToString() string {
	return fmt.Sprintf("Event Type %v", event.Type)
}
