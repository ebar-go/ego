package structure

func InArray[T comparable](arr []T, key T) bool {
	for _, t := range arr {
		if t == key {
			return true
		}
	}
	return false
}
