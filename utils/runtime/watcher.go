package runtime

import (
	"sync"
)

// Watcher
type Watcher interface {
	Stop()
}

// ChanWatcher implements the Watcher interface with channel.
type ChanWatcher struct {
	once   sync.Once
	stopCh chan struct{}
}

func (w *ChanWatcher) Stop() {
	w.once.Do(func() {
		close(w.stopCh)
	})
}

func (w *ChanWatcher) watch(chs ...chan struct{}) {
	defer HandleCrash()
	<-w.stopCh
	for _, ch := range chs {
		close(ch)
	}
}

func NewWatcher(chs ...chan struct{}) *ChanWatcher {
	watcher := &ChanWatcher{
		stopCh: make(chan struct{}),
	}
	go watcher.watch(chs...)
	return watcher
}
