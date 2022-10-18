package pool

type Worker interface {
	Schedule(fn func())
}
type WorkerPool struct {
	size int
	task chan struct{}
}

func NewWorkerPool(size int) *WorkerPool {
	pool := &WorkerPool{size: size, task: make(chan struct{}, size)}
	return pool
}

// Submit run some task
func (pool WorkerPool) Schedule(fn func()) {
	pool.task <- struct{}{}
	go func() {
		fn()
		<-pool.task
	}()
}

type GoroutinePool struct {
	work chan func()
}

func NewGoroutinePool(size int) *GoroutinePool {
	gp := &GoroutinePool{
		work: make(chan func(), size),
	}

	for i := 0; i < size; i++ {
		go gp.run()
	}
	return gp
}
func (p *GoroutinePool) Schedule(task func()) {
	select {
	case p.work <- task:
	default:
	}
}

func (p *GoroutinePool) run() {
	var task func()
	for {
		task = <-p.work
		task()
	}
}
