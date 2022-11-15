package pool

import (
	"sync"
	"sync/atomic"
)

type GoroutinePool struct {
	options Options
	// capacity of the pool.
	capacity int32

	// running is the number of the currently running goroutines.
	running int32

	// freeSignal is used to notice pool there are available
	// workers which can be sent to work.
	freeSignal chan struct{}

	// workers is a slice that store the available workers.
	workers []*Worker

	// stopped is used to check pool running status
	stopped int32

	// lock for synchronous operation
	lock sync.Mutex

	once sync.Once
}

type Options struct {
	Max  int32
	Idle int
	Step int
}

type Option func(options *Options)

func NewGoroutinePool(opts ...Option) *GoroutinePool {
	defaultOptions := Options{Max: 100, Idle: 10, Step: 10}
	for _, setter := range opts {
		setter(&defaultOptions)
	}
	pool := &GoroutinePool{
		options:    defaultOptions,
		capacity:   defaultOptions.Max,
		running:    0,
		freeSignal: make(chan struct{}, defaultOptions.Max),
		workers:    make([]*Worker, 0, 20),
	}

	// 提前准备空闲协程池
	pool.grow(defaultOptions.Idle)
	return pool
}

// grow 自动扩容worker数量
func (p *GoroutinePool) grow(n int) {
	for i := 0; i < n; i++ {
		// create instance of Worker
		w := NewWorker(p, 10)
		w.run()

		// push into the slice
		p.workers = append(p.workers, w)
		p.running++
		p.freeSignal <- struct{}{}
	}
}

type f func()

// Stop 停止协程池
func (p *GoroutinePool) Stop() {
	atomic.StoreInt32(&p.stopped, 1)
}

// Schedule 执行任务
func (p *GoroutinePool) Schedule(task func()) {

	// 判断pool是否已关闭
	if atomic.LoadInt32(&p.stopped) == 1 {
		return
	}

	// 调度一个worker来执行任务
	p.acquireWorker().submit(task)
}

// acquireWorker 获取一个worker实例
func (p *GoroutinePool) acquireWorker() (w *Worker) {
	p.lock.Lock()
	defer p.lock.Unlock()
	// 查看当前可用worker
	available := len(p.workers)
	if p.running < p.capacity {
		if available < 0 {
			p.grow(p.options.Idle)
		}

		<-p.freeSignal
		w = p.popLastWorker()
		return
	}

	// 当可用worker数为0且协程数达到上限时，
	// 因为此时已被lock住，且无法通过releaseWorker释放，所以会导致死锁
	// 所以这种情况下必须先释放锁
	p.lock.Unlock()
	<-p.freeSignal
	p.lock.Lock()
	w = p.popLastWorker()
	return w
}

func (p *GoroutinePool) popLastWorker() (w *Worker) {
	// 取数组最后一个worker
	n := len(p.workers) - 1
	w = p.workers[n]
	p.workers[n] = nil
	p.workers = p.workers[:n]
	return
}

// releaseWorker puts a worker back into free pool, recycling the goroutines.
func (p *GoroutinePool) releaseWorker(worker *Worker) {
	p.lock.Lock()
	p.workers = append(p.workers, worker)
	p.lock.Unlock()
	p.freeSignal <- struct{}{}
}
