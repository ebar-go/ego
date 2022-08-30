package component

import (
	"fmt"
	"testing"
)

func TestEventDispatcher(t *testing.T) {
	dispatcher := Provider().EventDispatcher()
	dispatcher.Listen("test", func(param any) {
		fmt.Println(param)
	})

	dispatcher.Trigger("test", 1)
}

func TestListenEvent(t *testing.T) {
	ListenEvent[string]("testString", func(s string) {
		fmt.Println(s)
	})

	ListenEvent[int]("testInt", func(n int) {
		fmt.Println(n)
	})

	Provider().EventDispatcher().Trigger("testString", "someString")
	Provider().EventDispatcher().Trigger("testInt", 1)
}

func TestEvent(t *testing.T) {
	e := Event[string]{Name: "test2"}
	e.Bind(func(param string) {
		fmt.Println(param)
	})
	Provider().EventDispatcher().Trigger("test2", "some other string")

	NewEvent[int]("test1").Bind(func(param int) {
		fmt.Println(param)
	})
	Provider().EventDispatcher().Trigger("test1", 1)

}
