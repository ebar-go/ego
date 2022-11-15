package pool

import (
	"sync"
)

type Pool struct {
	// capacity of the pool.
	capacity int32

	// running is the number of the currently running goroutines.
	running int32

	// freeSignal is used to notice pool there are available
	// workers which can be sent to work.
	freeSignal chan struct{}

	// workers is a slice that store the available workers.
	workers []*Worker

	// release is used to notice the pool to closed itself.
	release chan struct{}

	// lock for synchronous operation
	lock sync.Mutex

	once sync.Once
}

func NewGoroutinePool(cap int32) *Pool {
	pool := &Pool{
		capacity:   cap,
		running:    0,
		freeSignal: make(chan struct{}, cap),
		workers:    make([]*Worker, 0, 20),
		release:    make(chan struct{}, 1),
	}
	return pool
}

type f func()

func (p *Pool) Stop() {
	p.release <- struct{}{}
}

func (p *Pool) Schedule(task func()) {
	if len(p.release) > 0 {
		return
	}

	p.acquireWorker().submit(task)
}

func (p *Pool) acquireWorker() (w *Worker) {
	waiting := false
	p.lock.Lock()
	workers := p.workers
	n := len(workers) - 1
	if n < 0 {
		if p.running >= p.capacity {
			waiting = true
		} else {
			p.running++
		}
	} else {
		<-p.freeSignal
		w = workers[n]
		workers[n] = nil
		p.workers = workers[:n]
	}

	p.lock.Unlock()

	if waiting {
		<-p.freeSignal
		p.lock.Lock()
		workers = p.workers
		l := len(workers) - 1
		w = workers[l]
		workers[l] = nil
		p.workers = workers[:l]
		p.lock.Unlock()
	} else if w == nil {
		w = NewWorker(p, 10)
		w.run()
	}

	return w
}

// releaseWorker puts a worker back into free pool, recycling the goroutines.
func (p *Pool) releaseWorker(worker *Worker) {
	p.lock.Lock()
	p.workers = append(p.workers, worker)
	p.lock.Unlock()
	p.freeSignal <- struct{}{}
}
