package wait

import (
	"context"
	"github.com/ebar-go/ego/utils/runtime"
	"time"
)

// Until loops until stop channel is closed, running f every period.
func Until(f func(), period time.Duration, stopCh <-chan struct{}) {
	t := time.NewTimer(period)
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		func() {
			defer runtime.HandleCrash()
			f()
		}()
		t.Reset(period)

		select {
		case <-stopCh:
			if !t.Stop() {
				<-t.C
			}
			return
		case <-t.C:
		}
	}
}

// UntilWithContext loops until context is done, running f every period.
func UntilWithContext(f func(), period time.Duration, ctx context.Context) {
	Until(f, period, ctx.Done())
}
