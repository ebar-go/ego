package timewheel

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	tw := New(time.Second, 8)
	tw.Start()
	defer tw.Stop()

	tw.AfterFunc(time.Second, func() {
		fmt.Println("start")
	})
	time.Sleep(time.Second * 5)
}
