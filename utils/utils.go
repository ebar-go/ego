package utils

import (
	"log"
	"net"
	"errors"
	"reflect"
	"runtime"
	"fmt"
)

// CheckError
func CheckError(msg string, err error) {
	if err != nil {
		log.Printf("%s Error: %v\n", msg, err)
	}
}

// FatalError program will exit when err not nil
func FatalError(msg string, err error) {
	if err != nil {
		log.Fatalf("%s Error: %v\n", msg, err)
	}
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

// GetLocalIp return local ip
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

// Debug print params
func Debug(params ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("[Trace]%s[%d]:%v \n", file, line, params)
	}
}

// Trace return code trace info
func Trace() []string {
	trace := []string{}
	for i := 0; i < 10; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			trace = append(trace, fmt.Sprintf("[Trace]%s[%d]: \n", file, line))
		}
	}

	return trace
}

// GetKindOf return the kind of data
func GetKindOf(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}