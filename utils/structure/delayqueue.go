package structure

import (
	"container/heap"
	"sync"
	"sync/atomic"
	"time"
)

type item struct {
	Value    interface{} // 元素
	Priority int64       // 优先级
	Index    int         // 索引
}

// priorityQueue 实现heap.Interface接口，利用堆排序实现的按优先级排序的队列
type priorityQueue []*item

// newPriorityQueue 返回指定容量的队列实例
func newPriorityQueue(capacity int) priorityQueue {
	return make(priorityQueue, 0, capacity)
}

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	// 根据优先级堆元素进行排序
	return pq[i].Priority < pq[j].Priority
}

func (pq priorityQueue) Swap(i, j int) {
	// 交换元素，同时更新index属性
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	c := cap(*pq)
	if n+1 > c { // 触发扩容机制
		newCap := c * 2
		if c > 1024 { // 参考golang的slice扩容算法，当<=1024时，按倍数扩容，超过后就按1.25倍扩容
			newCap = c + c/4
		}
		npq := make(priorityQueue, n, newCap)
		copy(npq, *pq)
		*pq = npq
	}
	*pq = (*pq)[0 : n+1]
	item := x.(*item)
	item.Index = n
	(*pq)[n] = item
}

func (pq *priorityQueue) Pop() interface{} {
	n := len(*pq)
	c := cap(*pq)
	if n < (c/2) && c > 32 { // 触发缩容机制
		npq := make(priorityQueue, n, c/2)
		copy(npq, *pq)
		*pq = npq
	}
	item := (*pq)[n-1]
	item.Index = -1
	*pq = (*pq)[0 : n-1]
	return item
}

// PeekAndShift 根据优先级获取队列里满足条件的第一个元素
func (pq *priorityQueue) PeekAndShift(max int64) (*item, int64) {
	if pq.Len() == 0 {
		return nil, 0
	}

	// 因为队列是经过堆排序的，所以第一个元素始终是优先级最高的
	item := (*pq)[0]
	if item.Priority > max {
		// 如果未满足条件，返回差值
		return nil, item.Priority - max
	}

	// 满足条件，移除第一个数据
	heap.Remove(pq, 0)

	return item, 0
}

// DelayQueue 延时队列
type DelayQueue struct {
	Data     chan interface{} // 通过channel的方式输出数据，保证有序性
	mu       sync.Mutex       // 并发锁
	pq       priorityQueue    // 以时间为维度的优先级队列，时间越小，优先级越高
	sleeping int32            // 是否沉睡
	wakeup   chan struct{}    // 唤醒通知
	done     chan struct{}    // 停止通知
}

func NewDelayQueue(size int) *DelayQueue {
	return &DelayQueue{
		Data:   make(chan interface{}),
		pq:     newPriorityQueue(size),
		wakeup: make(chan struct{}),
		done:   make(chan struct{}),
	}
}

// Offer 插入一个数据以及设置其出列的时间
func (dq *DelayQueue) Offer(elem interface{}, expiration int64) {
	// 将时间设置为优先级，时间越小，优先级自然也就越高,也就会先出列
	item := &item{Value: elem, Priority: expiration}
	// 加锁，保证冰法安全
	dq.mu.Lock()
	// 通过堆实现插入，插入后的数组是有序的
	heap.Push(&dq.pq, item)
	index := item.Index
	dq.mu.Unlock()

	if index == 0 { // 判断到插入元素是队列的第一个元素时，则需要唤醒队列
		if atomic.CompareAndSwapInt32(&dq.sleeping, 1, 0) {
			dq.wakeup <- struct{}{}
		}
	}
}

// Poll 根据时间函数读取队列数据
func (dq *DelayQueue) Poll(nowF func() int64) {
	for {
		now := nowF()

		dq.mu.Lock()
		// 尝试从队列里获取一个数据
		item, delta := dq.pq.PeekAndShift(now)
		if item == nil {
			// 代表队列里暂无数据可读取到,将队列设置为沉睡状态
			atomic.StoreInt32(&dq.sleeping, 1)
		}
		dq.mu.Unlock()

		if item == nil {
			if delta == 0 { // 代表队列里没有任何元素
				select {
				case <-dq.wakeup: // 通过wakeup阻塞，等待新的元素插入再运行
					continue
				case <-dq.done:
					goto exit
				}
			} else if delta > 0 { // 代表队列里至少有一个元素，只是未到时间
				select {
				case <-dq.wakeup: // 如果此刻有新元素插入到队列里，可能是优先级更高，所以需要循环读取
					continue
				case <-time.After(time.Duration(delta) * time.Millisecond):
					// 通过timer实现等待，此刻已满足读取队列中第一个元素的条件
					if atomic.SwapInt32(&dq.sleeping, 0) == 0 {
						// A caller of Offer() is being blocked on sending to wakeupC,
						// drain wakeupC to unblock the caller.
						<-dq.wakeup
					}
					continue
				case <-dq.done:
					goto exit
				}
			}
		}

		// 从队列里成功读取到数据，将数据写入到data的channel里，供调用方读取
		select {
		case dq.Data <- item.Value:
		case <-dq.done:
			goto exit
		}
	}

exit:
	// Reset the states
	atomic.StoreInt32(&dq.sleeping, 0)
}

// Close 关闭
func (dq *DelayQueue) Close() {
	close(dq.done)
	close(dq.Data)
}
