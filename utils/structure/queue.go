package structure

// Queue represents a queue of any object
type Queue[T any] struct {
	list chan T
}

// Offer push the active fd to the queue, it will be deposed when queue is full
func (queue *Queue[T]) Offer(items ...T) {
	for _, item := range items {
		// depose fd when queue is full
		select {
		case queue.list <- item:
		default:
			// must add default option,
			// otherwise it will:
			//		fatal error: all goroutines are asleep- deadlock!

		}
	}
}

// Polling poll with callback function
func (queue *Queue[T]) Polling(stopCh <-chan struct{}, handler func(item T)) {
	for {
		select {
		// stop when signal is closed
		case <-stopCh:
			return
		case item := <-queue.list:
			handler(item)
		}
	}
}

// Empty returns true if the queue is empty
func (queue *Queue[T]) Empty() bool {
	return queue.Length() == 0
}

// Length returns the length of queue
func (queue *Queue[T]) Length() int {
	return len(queue.list)
}

// NewQueue returns a new queue with the given size
func NewQueue[T any](size int) *Queue[T] {
	return &Queue[T]{
		list: make(chan T, size),
	}
}
