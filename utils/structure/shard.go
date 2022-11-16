package structure

type ShardContainer[T any] []T

func NewShardContainer[T any](size int, constructor func() T) ShardContainer[T] {
	sc := make(ShardContainer[T], size)
	for i := 0; i < size; i++ {
		sc[i] = constructor()
	}
	return sc
}

func (shard ShardContainer[T]) GetShard(num int) T {
	idx := num & (len(shard) - 1)
	return shard[idx]
}

func (shard ShardContainer[T]) Iterator(iterator func(T)) {
	for _, item := range shard {
		iterator(item)
	}
}
