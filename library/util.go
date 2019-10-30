package library

import (
	"math"
	"runtime"
	"fmt"
)

func Round(f float64) int {
	return int(math.Floor(f + 0.5))
}

// Debug
func Debug(params ...interface{})  {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("[Trace]%s[%d]:\n", file, line)
		fmt.Println(params)
	}
}

// ArrayUniqueInt
func ArrayUniqueInt(items []int) []int {
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

// GetStringItemOfArray
func GetStringItemOfArray(items map[int]string, index int, defaultValue string) string {
	if _, ok := items[index]; ok {
		return items[index]
	}

	return defaultValue
}

func Float64ToString(a float64) string {
	return fmt.Sprintf("%.f", a)
}