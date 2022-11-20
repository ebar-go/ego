package structure

import "github.com/ebar-go/ego/utils/number"

type Sharding[T any] []T

func NewSharding[T any](size int, constructor func() T) Sharding[T] {
	size = number.RoundUp(size)
	sc := make(Sharding[T], size)
	for i := 0; i < size; i++ {
		sc[i] = constructor()
	}
	return sc
}

func (shard Sharding[T]) GetShard(num int) T {
	idx := num & (len(shard) - 1)
	return shard[idx]
}

func (shard Sharding[T]) Iterator(iterator func(T)) {
	for _, item := range shard {
		iterator(item)
	}
}
