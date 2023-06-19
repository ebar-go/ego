package pool

import (
	"sync"
	"sync/atomic"
	"time"
)

type GoroutinePool interface {
	// Schedule 调度worker执行任务
	Schedule(task func())

	// Stop 停止工作
	Stop()
}

type ScalableGroutinePool struct {
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
	done chan struct{}
}

type Options struct {
	Max     int32 // maximum number of workers
	Idle    int   // idle number of workers
	Block   bool
	Timeout time.Duration // quit time for worker
}

type Option func(options *Options)

func NewGoroutinePool(opts ...Option) *ScalableGroutinePool {
	defaultOptions := Options{Max: 100, Idle: 10, Timeout: time.Second * 10, Block: true}
	for _, setter := range opts {
		setter(&defaultOptions)
	}
	pool := &ScalableGroutinePool{
		options:    defaultOptions,
		capacity:   defaultOptions.Max,
		running:    0,
		freeSignal: make(chan struct{}, defaultOptions.Max),
		workers:    make([]*Worker, 0, 1024),
		done:       make(chan struct{}),
	}

	// 提前准备空闲协程池
	pool.grow(defaultOptions.Idle)
	go pool.monitor()

	return pool
}

func (p *ScalableGroutinePool) monitor() {
	for {
		select {
		case <-p.done:
			return
		default:
			// 定时缩容
			p.scaleDown()
		}
	}
}

// grow 自动扩容worker数量
func (p *ScalableGroutinePool) grow(n int) {
	for i := 0; i < n; i++ {
		// create instance of Worker
		w := NewWorker(10, func() {
			atomic.AddInt32(&p.running, -1)
		}, p.releaseWorker)

		// push into the slice
		p.workers = append(p.workers, w)
		p.running++
		p.freeSignal <- struct{}{}
	}
}

type f func()

// Stop 停止协程池
func (p *ScalableGroutinePool) Stop() {
	p.once.Do(func() {
		atomic.StoreInt32(&p.stopped, 1)
		close(p.done)
		p.lock.Lock()
		for _, worker := range p.workers {
			worker.Stop()
		}
		p.lock.Unlock()
	})
}

// Schedule 执行任务
func (p *ScalableGroutinePool) Schedule(task func()) {
	// 判断pool是否已关闭
	if atomic.LoadInt32(&p.stopped) == 1 {
		return
	}

	// 调度一个worker来执行任务
	p.acquireWorker().Submit(task, p.options.Block)
}

// acquireWorker 获取一个worker实例
func (p *ScalableGroutinePool) acquireWorker() (w *Worker) {
	p.lock.Lock()
	defer p.lock.Unlock()
	// 查看当前可用worker
	available := len(p.workers)
	if p.running < p.capacity {
		if available == 0 {
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

// popLastWorker 获取空闲队列里最尾部的一个worker
func (p *ScalableGroutinePool) popLastWorker() (w *Worker) {
	// 取数组最后一个worker
	n := len(p.workers) - 1
	w = p.workers[n]
	p.workers[n] = nil
	p.workers = p.workers[:n]
	return
}

// releaseWorker puts a worker back into free pool, recycling the goroutines.
func (p *ScalableGroutinePool) releaseWorker(worker *Worker) {
	p.lock.Lock()
	p.workers = append(p.workers, worker)
	p.lock.Unlock()
	p.freeSignal <- struct{}{}
}

// scaleDown 缩容
func (p *ScalableGroutinePool) scaleDown() {
	// 根据时间控制每隔一段时间按策略缩容
	time.Sleep(p.options.Timeout)
	p.lock.Lock()
	defer p.lock.Unlock()

	// 低于空闲数量则不缩容
	available := len(p.workers)
	if available <= p.options.Idle {
		return
	}

	num := (available - p.options.Idle) / 4
	for i := 0; i < num; i++ {
		p.workers[i].Stop()
	}
	p.workers = p.workers[num:]
}

// FixedGoroutinePool 固定数量的协程池
type FixedGoroutinePool struct {
	task  chan f
	block bool // 控制当任务已满时是否阻塞

	once sync.Once
	done chan struct{}
}

func (pool *FixedGoroutinePool) Schedule(task func()) {
	if pool.block {
		// 阻塞等待任务被执行
		pool.task <- task
		return
	}

	// 非阻塞模式
	select {
	case pool.task <- task:
	default:
	}

}

func (pool *FixedGoroutinePool) Stop() {
	pool.once.Do(func() {
		close(pool.done)
	})
}

func (pool *FixedGoroutinePool) run() {
	for {
		select {
		case <-pool.done:
			return
		case fn := <-pool.task:
			fn()
		}
	}
}

type FixedOptions struct {
	Num   int  // 协程数
	Block bool // 任务已满时是否需要阻塞
}
type fixedOption func(options *FixedOptions)

func NewFixedGoroutinePool(options ...fixedOption) *FixedGoroutinePool {
	defaultOptions := FixedOptions{Num: 100, Block: true}
	for _, setter := range options {
		setter(&defaultOptions)
	}
	pool := &FixedGoroutinePool{
		task:  make(chan f, defaultOptions.Num),
		block: defaultOptions.Block,
		done:  make(chan struct{}),
	}

	// 直接启动固定数量的协程
	for i := 0; i < defaultOptions.Num; i++ {
		go pool.run()
	}

	return pool
}
