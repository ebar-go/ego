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

func BenchmarkTimer(b *testing.B) {
	// cpu is cost fast when create so much timer
	// time-wheel timer is better
	n := 400000
	b.Run("timewheel", func(b *testing.B) {
		b.ReportAllocs()
		tw := New(time.Second, 64)
		tw.Start()
		defer tw.Stop()
		for i := 0; i < n; i++ {
			tw.AfterFunc(time.Second, func() {
				_ = fmt.Sprintf("benchmark timer")
			})
		}
	})
	b.Run("normal", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < n; i++ {
			time.AfterFunc(time.Second, func() {
				_ = fmt.Sprintf("benchmark timer")
			})
		}
	})
}
