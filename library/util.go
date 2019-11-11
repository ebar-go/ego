package library

import (
	"math"
	"runtime"
	"fmt"
	"os"
	"encoding/json"
	"net"
	"errors"
	"strings"
	"strconv"
)

// JsonEncode json序列号
func JsonEncode(v interface{}) (string, error) {
	bytes , err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// GetCurrentPath 获取当前路径
func GetCurrentPath() string {
	path, _ := os.Getwd()
	return path
}

// Round 四舍五入取整
func Round(f float64) int {
	return int(math.Floor(f + 0.5))
}

// Debug 打印信息
func Debug(params ...interface{})  {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("[Trace]%s[%d]:\n%v \n", file, line, params)
	}
}

// Trace 返回trace日志
func Trace() []string {
	trace := []string{}
	for i:=0; i<5; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			trace = append(trace, fmt.Sprintf("[Trace]%s[%d]: \n", file, line))
		}
	}

	return trace
}

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

// Float64ToString float转string
func Float64ToString(a float64) string {
	return fmt.Sprintf("%.f", a)
}

// GetLocalIp 获取本地IP
func GetLocalIp() (string, error) {
	addressItems, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addressItems {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}

	return "", errors.New("failed to get local address")
}

// Implode 连接slice为字符串
func Implode(separator string ,items []interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(items), "[]"), " ", separator, -1)
}

// Explode 分割字符串为int
func ExplodeInt(str, separator string) (result []int) {
	var strItems = strings.Split(str, separator)
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
func ExplodeString(str, separator string) (result []string) {

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

// MergeMaps 合并
func MergeMaps(items ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, item := range items {
		for key, value := range item {
			result[key] = value
		}
	}

	return result
}