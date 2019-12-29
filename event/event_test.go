package event

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	dispatcher := NewDispatcher()
	listener := NewListener(func(event Event) {
		fmt.Println(event.Type, event.Dispatcher, event.Params)
	})

	dispatcher.AddListener("TEST", listener)

	dispatcher.DispatchEvent(New("TEST", map[string]interface{}{"hello":"world"}))
}
