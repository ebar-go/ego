package timewheel

import (
	"container/list"
	"errors"
	"github.com/ebar-go/ego/utils/structure"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type waitGroupWrapper struct {
	sync.WaitGroup
}

func (w *waitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}

// TimeWheel 实现一个时间轮定时器，在指定时间运行指定的任务
type TimeWheel struct {
	tick int64 // 时间跨度，单位是ms
	size int64 // 时间轮格数

	interval    int64 // 总体时间跨度，interval=tick × size
	currentTime int64 // 当前运行的时间,是 tickMs 的整数倍

	buckets    []*bucket // 存储任务的桶
	stop       chan struct{}
	waitGroup  waitGroupWrapper
	delayQueue *structure.DelayQueue

	// 层级时间轮，类似于链表结构，每一层时间轮都会指向下一层时间轮
	overflowWheel unsafe.Pointer
}

// New 实例化一个时间轮定时器
func New(tick time.Duration, wheelSize int64) *TimeWheel {
	// 根据tick计算定时器的执行间隔，单位是毫秒
	tickMs := int64(tick / time.Millisecond)
	if tickMs <= 0 {
		panic(errors.New("tick must be greater than or equal to 1ms"))
	}

	// 初始化一个延迟队列，任何层的时间轮都是共用一个队列
	delayQueue := structure.NewDelayQueue(int(wheelSize))
	return newTimeWheel(tickMs, wheelSize, UnixMilli(time.Now()), delayQueue)
}

func newTimeWheel(tickMs int64, wheelSize int64, startMs int64, delayQueue *structure.DelayQueue) *TimeWheel {
	// 初始化任务桶
	buckets := make([]*bucket, wheelSize)
	for i := 0; i < int(wheelSize); i++ {
		buckets[i] = newBucket()
	}

	return &TimeWheel{
		tick:        tickMs,
		size:        wheelSize,
		interval:    tickMs * wheelSize,
		currentTime: truncate(startMs, tickMs),
		buckets:     buckets,
		stop:        make(chan struct{}),
		delayQueue:  delayQueue,
	}
}

// advanceClock 移动指针
func (tw *TimeWheel) advanceClock(expiration int64) {
	// 加载当前时间
	currentTime := atomic.LoadInt64(&tw.currentTime)
	if expiration >= currentTime+tw.tick {
		// 将时间前进一个时间刻度
		currentTime = truncate(expiration, tw.tick)
		// 更新currentTime
		atomic.StoreInt64(&tw.currentTime, currentTime)

		// 同时将更高层级的时间轮的currentTime也前拨
		overflowWheel := atomic.LoadPointer(&tw.overflowWheel)
		if overflowWheel != nil {
			(*TimeWheel)(overflowWheel).advanceClock(currentTime)
		}
	}
}

func (tw *TimeWheel) Start() {
	// 开启协程，执行拉取任务的逻辑
	tw.waitGroup.Wrap(func() {
		tw.delayQueue.Poll(func() int64 {
			return UnixMilli(time.Now())
		})
	})

	// 开启协程，执行执行任务的逻辑
	tw.waitGroup.Wrap(func() {
		for {
			select {
			case item := <-tw.delayQueue.Data:
				if item == nil {
					continue
				}
				b := item.(*bucket)
				// 调整指针
				tw.advanceClock(b.expiration)
				// 对任务列表进行遍历，将到期的任务执行，未到期的继续插入
				b.Flush(tw.addOrRun)
			case <-tw.stop:
				return

			}
		}
	})
}

// Stop 关闭
func (tw *TimeWheel) Stop() {
	close(tw.stop)
	tw.delayQueue.Close()
	tw.waitGroup.Wait()
}

// AfterFunc
func (tw *TimeWheel) AfterFunc(d time.Duration, f func()) *Timer {
	t := &Timer{
		expiration: UnixMilli(time.Now().Add(d)),
		task:       f,
	}

	tw.addOrRun(t)
	return t
}

// add 添加任务，如果任务已过期，则返回false
func (tw *TimeWheel) add(t *Timer) bool {
	// 加载当前时间
	currentTime := atomic.LoadInt64(&tw.currentTime)
	// 判断任务的执行时间是否在当前时间轮的执行期内
	if t.expiration < currentTime+tw.tick {
		// 返回false代表让任务直接执行
		return false
	} else if t.expiration < currentTime+tw.interval {
		// 没到执行时间且在本层时间轮里执行
		// 找到对应的bucket
		virtualID := t.expiration / tw.tick
		b := tw.buckets[virtualID%tw.size]
		b.Add(t)

		// 更新执行时间成功后，需要调整队列的优先级,如果时间相同，则不需要调整
		if b.SetExpiration(virtualID * tw.tick) {
			tw.delayQueue.Offer(b, b.expiration)
		}

		return true
	} else { // 当执行时间在此周期外
		overflowWheel := atomic.LoadPointer(&tw.overflowWheel)
		if overflowWheel == nil {
			// 初始化下一层的时间轮
			atomic.CompareAndSwapPointer(
				&tw.overflowWheel,
				nil,
				unsafe.Pointer(newTimeWheel(tw.interval, tw.size, currentTime, tw.delayQueue)),
			)
			// 再加载一次
			overflowWheel = atomic.LoadPointer(&tw.overflowWheel)
		}
		// 利用递归的思想，将任务插入到对应层的bucket里
		return (*TimeWheel)(overflowWheel).add(t)
	}

}
func (tw *TimeWheel) addOrRun(t *Timer) {
	if !tw.add(t) {
		// 任务已过期，直接执行任务
		go t.task()
	}
}

// truncate returns the result of rounding x toward zero to a multiple of m.
// If m <= 0, Truncate returns x unchanged.
func truncate(x, m int64) int64 {
	if m <= 0 {
		return x
	}
	return x - x%m
}

// bucket 存储任务列表的桶
type bucket struct {
	expiration int64 // 过期时间,需要用atomic保证并发安全

	mu     sync.Mutex
	timers *list.List // 底层采用的是链表结构
}

func (b *bucket) Expiration() int64 {
	return atomic.LoadInt64(&b.expiration)
}

func (b *bucket) SetExpiration(expiration int64) bool {
	return atomic.SwapInt64(&b.expiration, expiration) != expiration
}

// Flush 清空
func (b *bucket) Flush(reinsert func(t *Timer)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	// 遍历链表
	for e := b.timers.Front(); e != nil; {
		next := e.Next()

		t := e.Value.(*Timer)
		b.remove(t)

		// 满足条件的直接执行，否则继续插入到其他bucket，等待到期再执行
		reinsert(t)

		e = next
	}

	b.SetExpiration(-1)

}

// Remove
func (b *bucket) Remove(t *Timer) bool {
	// 需要保证线程安全
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.remove(t)
}

// remove 移除队列里的元素
func (b *bucket) remove(t *Timer) bool {
	if t.getBucket() != b {
		return false
	}
	b.timers.Remove(t.element)
	t.setBucket(nil)
	t.element = nil
	return true
}

// Add 添加任务列表
func (b *bucket) Add(t *Timer) {
	b.mu.Lock()
	elem := b.timers.PushBack(t)
	t.setBucket(b)
	t.element = elem
	b.mu.Unlock()

}

func newBucket() *bucket {
	return &bucket{
		expiration: -1,
		timers:     list.New(),
	}
}

type Timer struct {
	expiration int64 // in milliseconds

	// 当出现并发读写的时候，对于指针都可以用unsafe.Pointer处理
	bucket  unsafe.Pointer // bucket,需要用atomic保证并发安全
	element *list.Element
	task    func() // 具体任务
}

func (t *Timer) getBucket() *bucket {
	return (*bucket)(atomic.LoadPointer(&t.bucket))
}

func (t *Timer) setBucket(bucket *bucket) {
	atomic.StorePointer(&t.bucket, unsafe.Pointer(bucket))
}

// Stop 停止任务,当任务已执行时返回false
func (t *Timer) Stop() (stopped bool) {
	b := t.getBucket()
	for b != nil {
		stopped = b.Remove(t)
		b = t.getBucket()
	}
	return
}

// Reset 将定时任务往后延时执行
func (t *Timer) Reset(d time.Duration) bool {
	now := time.Now()
	if t.expiration <= UnixMilli(now) { // 说明任务已执行完成
		return false
	}
	t.expiration = UnixMilli(now.Add(d))
	return true
}

// UnixMilli 毫秒
func UnixMilli(t time.Time) int64 {
	return t.UnixNano() / 1e6
}
