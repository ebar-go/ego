package number

import "strconv"

// RoundUp 返回邻近的2的N次方的数
func RoundUp(v int) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}

func Int(n interface{}) int {
	switch v := n.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		num, _ := strconv.Atoi(v)
		return num

	}
	return 0
}
