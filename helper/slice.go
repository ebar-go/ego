package helper

import (
	"fmt"
	"strconv"
	"strings"
)

// ArrayUniqueInt 排重
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

// SliceIntToString 切片int转字符串
func SliceIntToString(items []int) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, strconv.Itoa(item))
	}
	return result
}

// GetStringItemOfArray 获取数组元素，没有就返回默认值
func GetStringItemOfArray(items map[int]string, index int, defaultValue string) string {
	if _, ok := items[index]; ok {
		return items[index]
	}

	return defaultValue
}

// Implode 连接slice为字符串
func Implode(separator string, items []interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(items), "[]"), " ", separator, -1)
}

// Explode 分割字符串为int
func ExplodeInt(str, separator string) (result []int) {
	strItems := strings.Split(str, separator)
	for _, v := range strItems {
		mid, e := strconv.Atoi(strings.Trim(v, ""))
		if e != nil {
			continue
		}
		result = append(result, mid)
	}
	return
}

// Explode 分割字符串为int
func ExplodeString(str, separator string) []string {
	return strings.Split(str, separator)
}

// IntSliceToInterface int类型的切片转interface
func IntSliceToInterface(items []int) []interface{} {
	var interfaceSlice []interface{} = make([]interface{}, len(items))
	for i, d := range items {
		interfaceSlice[i] = d
	}

	return interfaceSlice
}
