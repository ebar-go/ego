package convert

func ToSlice[T any](items []interface{}) []T {
	res := make([]T, 0)
	for _, item := range items {
		res = append(res, item.(T))
	}
	return res
}
