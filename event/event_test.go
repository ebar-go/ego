package event

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	eventType := "HELLO"
	eventParams := "world"
	ev := New(eventType, eventParams)
	assert.Equal(t, eventType, ev.Type)
	assert.Equal(t, eventParams, ev.Params)
}



func TestNewDispatcher(t *testing.T) {
	assert.NotNil(t, NewDispatcher())
}

func TestNewListener(t *testing.T) {
	handler := func(ev Event) {
		fmt.Println(ev.Type, ev.Dispatcher, ev.Params)
	}
	listener := NewListener(handler)
	assert.NotNil(t, listener)
}

func TestEvent_Clone(t *testing.T) {
	eventType := "HELLO"
	eventParams := "world"
	ev := New(eventType, eventParams)

	eventClone := ev.Clone()
	assert.NotNil(t, eventClone)
	assert.Equal(t, eventType, eventClone.Type)
	assert.Equal(t, eventParams, eventClone.Params)
}

func TestEvent_ToString(t *testing.T) {
	ev := New("HELLO", "world")
	fmt.Println(ev.ToString())
}

func TestEventDispatcher_AddListener(t *testing.T) {
	dispatcher := NewDispatcher()
	listener := NewListener(func(ev Event) {
		fmt.Println(ev.Type, ev.Dispatcher, ev.Params)
	})

	dispatcher.AddListener("TEST", listener)
	assert.True(t, dispatcher.HasListener("TEST"))

}

func TestEventDispatcher_RemoveListener(t *testing.T) {
	dispatcher := NewDispatcher()
	listener := NewListener(func(ev Event) {
		fmt.Println(ev.Type, ev.Dispatcher, ev.Params)
	})

	dispatcher.AddListener("TEST", listener)
	assert.True(t, dispatcher.HasListener("TEST"))
	res := dispatcher.RemoveListener("TEST", listener)
	assert.True(t, res)

}

func TestEventDispatcher_HasListener(t *testing.T) {
	dispatcher := NewDispatcher()
	listener := NewListener(func(ev Event) {
		fmt.Println(ev.Type, ev.Dispatcher, ev.Params)
	})

	dispatcher.AddListener("TEST", listener)
	assert.True(t, dispatcher.HasListener("TEST"))
	assert.False(t, dispatcher.HasListener("NOT_EXIST"))

}

func TestEventDispatcher_DispatchEvent(t *testing.T) {
	dispatcher := NewDispatcher()
	listener := NewListener(func(ev Event) {
		fmt.Println(ev.Type, ev.Dispatcher, ev.Params)
	})

	dispatcher.AddListener("TEST", listener)

	dispatcher.DispatchEvent(New("TEST", map[string]interface{}{"hello":"world"}))
}
