package array

import "strconv"

// UniqueInt remove duplicate elements
func UniqueInt(items []int) []int {
	result := make([]int, 0, len(items))
	temp := map[int]struct{}{}
	for _, item := range items {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// Int2String return strings slice
func Int2String(items []int) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, strconv.Itoa(item))
	}
	return result
}

// Int2Interface return interface slice
func Int2Interface(items []int) []interface{} {
	var interfaceSlice = make([]interface{}, len(items))
	for i, d := range items {
		interfaceSlice[i] = d
	}

	return interfaceSlice
}