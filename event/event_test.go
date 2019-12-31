package event

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	eventType := "HELLO"
	eventParams := "world"
	event := New(eventType, eventParams)
	assert.Equal(t, eventType, event.Type)
	assert.Equal(t, eventParams, event.Params)
}



func TestNewDispatcher(t *testing.T) {
	assert.NotNil(t, NewDispatcher())
}

func TestNewListener(t *testing.T) {
	handler := func(event Event) {
		fmt.Println(event.Type, event.Dispatcher, event.Params)
	}
	listener := NewListener(handler)
	assert.NotNil(t, listener)
}

func TestEvent_Clone(t *testing.T) {
	eventType := "HELLO"
	eventParams := "world"
	event := New(eventType, eventParams)

	eventClone := event.Clone()
	assert.NotNil(t, eventClone)
	assert.Equal(t, eventType, eventClone.Type)
	assert.Equal(t, eventParams, eventClone.Params)
}

func TestEvent_ToString(t *testing.T) {
	event := New("HELLO", "world")
	fmt.Println(event.ToString())
}

func TestEventDispatcher_AddListener(t *testing.T) {
	dispatcher := NewDispatcher()
	listener := NewListener(func(event Event) {
		fmt.Println(event.Type, event.Dispatcher, event.Params)
	})

	dispatcher.AddListener("TEST", listener)
	assert.True(t, dispatcher.HasListener("TEST"))

}

func TestEventDispatcher_RemoveListener(t *testing.T) {
	dispatcher := NewDispatcher()
	listener := NewListener(func(event Event) {
		fmt.Println(event.Type, event.Dispatcher, event.Params)
	})

	dispatcher.AddListener("TEST", listener)
	assert.True(t, dispatcher.HasListener("TEST"))
	res := dispatcher.RemoveListener("TEST", listener)
	assert.True(t, res)

}

func TestEventDispatcher_HasListener(t *testing.T) {
	dispatcher := NewDispatcher()
	listener := NewListener(func(event Event) {
		fmt.Println(event.Type, event.Dispatcher, event.Params)
	})

	dispatcher.AddListener("TEST", listener)
	assert.True(t, dispatcher.HasListener("TEST"))
	assert.False(t, dispatcher.HasListener("NOT_EXIST"))

}

func TestEventDispatcher_DispatchEvent(t *testing.T) {
	dispatcher := NewDispatcher()
	listener := NewListener(func(event Event) {
		fmt.Println(event.Type, event.Dispatcher, event.Params)
	})

	dispatcher.AddListener("TEST", listener)

	dispatcher.DispatchEvent(New("TEST", map[string]interface{}{"hello":"world"}))
}
