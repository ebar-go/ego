package array

import (
	"strconv"
)

// StringSlice
type StringSlice struct {
	items []string
}

// Int return IntSlice
func String(items []string) StringSlice {
	return StringSlice{items: items}
}

// Length
func (s StringSlice) Length() int {
	return len(s.items)
}

// Push
func (s *StringSlice) Push(items ...string) {
	s.items = append(s.items, items...)
}

// Unique return unique elem slice
func (s StringSlice) Unique() []string {
	result := make([]string, 0, len(s.items))
	temp := map[string]struct{}{}
	for _, item := range s.items {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// ToString return string slice
func (s StringSlice) ToInt() []int {
	result := make([]int, 0, len(s.items))
	for _, item := range s.items {
		i, err := strconv.Atoi(item)
		if err != nil {
			continue
		}
		result = append(result, i)
	}
	return result
}

// Implode
func (s StringSlice) Implode(separator string) string {
	return Implode(s.items, separator)
}

// Items
func (s StringSlice) Items() []string {
	return s.items
}

// Has
func (s StringSlice) Has(elem string) bool {
	for _, i := range s.items {
		if elem == i {
			return true
		}
	}

	return false
}
