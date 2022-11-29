package pool

import (
	"sync"
)

type Worker struct {
	task chan f
	once sync.Once
	done chan struct{}
}

func NewWorker(cap int, beforeCloseCallback func(), afterCallback func(w *Worker)) *Worker {
	w := &Worker{
		task: make(chan f, cap),
		done: make(chan struct{}),
	}
	go w.run(beforeCloseCallback, afterCallback)
	return w
}

func (w *Worker) run(beforeCloseCallback func(), afterCallback func(w *Worker)) {
	var (
		fn     f
		opened bool
	)
	for {
		select {
		case <-w.done:
			beforeCloseCallback()
			return
		case fn, opened = <-w.task:
			if !opened {
				return
			}
			fn()
			//回收复用
			afterCallback(w)
		}
	}
}

// stop this worker.
func (w *Worker) Stop() {
	w.once.Do(func() {
		close(w.done)
		close(w.task)
	})
}

// Submit sends a task to this worker.
func (w *Worker) Submit(task f, block bool) {
	if block {
		select {
		case <-w.done:
		case w.task <- task:
		}
		return
	}

	select {
	case w.task <- task:
	default:
	}
}
