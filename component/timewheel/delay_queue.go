package timewheel

import (
	"container/heap"
	"sync"
	"sync/atomic"
	"time"
)

// DelayQueue 延时队列
type DelayQueue struct {
	C        chan interface{} // 利用channel来传输任务
	wakeupC  chan struct{}    // 有新任务就唤醒队列
	mu       sync.Mutex
	pq       priorityQueue // 延时队列
	sleeping int32         // 是否为睡眠状态
}

func NewDelayQueue(size int) *DelayQueue {
	return &DelayQueue{
		C:       make(chan interface{}),
		pq:      newPriorityQueue(size),
		wakeupC: make(chan struct{}),
	}
}

type Item struct {
	value    interface{}
	priority int64 // 优先级
	index    int
}

// Channel 返回
func (queue *DelayQueue) Channel() chan interface{} {
	return queue.C
}

// Poll 拉取任务
func (queue *DelayQueue) Poll(exitC chan struct{}, nowF func() int64) {
	for {
		now := nowF()
		queue.mu.Lock()

		// 根据当前时间，获取一个新的任务
		// delta是距离任务执行的时间差额
		item, delta := queue.pq.peekAndShift(now)
		if item == nil {
			// 没有任务，进入睡眠状态
			atomic.StoreInt32(&queue.sleeping, 1)
		}

		queue.mu.Unlock()
		if item == nil { // 没有立即执行的任务
			if delta == 0 { // 没任务
				select {
				case <-queue.wakeupC: // 利用channel实现等待，直到有新任务才继续执行
					continue
				case <-exitC:
					goto exit
				}
			} else if delta > 0 { // 有个即将执行的任务
				select {
				case <-queue.wakeupC: // 同上
					continue
				case <-time.After(time.Duration(delta) * time.Millisecond): // 时间到就执行
					if atomic.SwapInt32(&queue.sleeping, 0) == 0 {
						<-queue.wakeupC
					}
					continue
				case <-exitC:
					goto exit
				}
			}
		}

		// 有任务满足执行的条件，将任务放入channel
		select {
		case queue.C <- item.value:
		case <-exitC:
			goto exit

		}
	}

exit:
	// Reset the states
	atomic.StoreInt32(&queue.sleeping, 0)
}

// Offer 往队列里加入一个任务
func (queue *DelayQueue) Offer(val interface{}, expiration int64) {
	// 把时间当做队列的优先级，时间越小，优先级越高，越先执行
	item := &Item{value: val, priority: expiration}

	queue.mu.Lock()
	heap.Push(&queue.pq, item) // 利用堆进行排序
	index := item.index
	queue.mu.Unlock()

	if index == 0 { // 经过堆排序后的index为0,代表加入的是一个优先级更高的任务
		if atomic.CompareAndSwapInt32(&queue.sleeping, 1, 0) {
			// 通知队列开始工作
			queue.wakeupC <- struct{}{}
		}
	}

}

// priorityQueue 优先级队列，利用堆结构进行排序
type priorityQueue []*Item

func newPriorityQueue(capacity int) priorityQueue {
	return make(priorityQueue, 0, capacity)
}

func (queue priorityQueue) Len() int {
	return len(queue)
}

func (queue priorityQueue) Less(i, j int) bool {
	return queue[i].priority < queue[j].priority
}

func (queue priorityQueue) Swap(i, j int) {
	queue[i], queue[j] = queue[j], queue[i]
	queue[i].index = i
	queue[j].index = j
}

func (queue *priorityQueue) Push(x any) {
	n := len(*queue)
	c := cap(*queue)
	if n+1 > c { // 成倍扩容数组
		temp := make(priorityQueue, n, c*2)
		copy(temp, *queue)
		*queue = temp
	}

	*queue = (*queue)[:n+1]
	item := x.(*Item)
	item.index = n
	(*queue)[n] = item
}

func (queue *priorityQueue) Pop() any {
	n := len(*queue)
	c := cap(*queue)
	if n < (c/2) && c > 25 { // 按倍缩容，且容量不能小于25
		temp := make(priorityQueue, n, c/2)
		copy(temp, *queue)
		*queue = temp
	}
	item := (*queue)[n-1]
	item.index = -1
	*queue = (*queue)[0 : n-1]
	return item
}

// PeekAndShift 取第一个
func (queue *priorityQueue) peekAndShift(max int64) (*Item, int64) {
	if queue.Len() == 0 {
		return nil, 0
	}

	// 取第一位
	item := (*queue)[0]
	if item.priority > max {
		// 如果任务的执行时间还未到，返回差值
		return nil, item.priority - max
	}

	// 时间已到，直接移除元素并返回
	heap.Remove(queue, 0)
	return item, 0
}
