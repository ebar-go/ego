package wait

import (
	"github.com/ebar-go/ego/utils/runtime"
	"time"
)

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
