package pool

import (
	"sync"
	"sync/atomic"
)

type Worker struct {
	task chan f
	pool *GoroutinePool
	once sync.Once
	done chan struct{}
}

func NewWorker(pool *GoroutinePool, size int) *Worker {
	return &Worker{
		pool: pool,
		task: make(chan f, size),
		done: make(chan struct{}),
	}
}

func (w *Worker) run() {
	go func() {
		for {
			select {
			case <-w.done:
				atomic.AddInt32(&w.pool.running, -1)
				return
			case fn := <-w.task:
				fn()
				//回收复用
				w.pool.releaseWorker(w)
			}
		}
	}()
}

// stop this worker.
func (w *Worker) stop() {
	w.once.Do(func() {
		close(w.done)
	})
}

// submit sends a task to this worker.
func (w *Worker) submit(task f) {
	select {
	case w.task <- task:
	default:
	}
}
