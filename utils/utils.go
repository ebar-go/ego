package utils

import (
	"errors"
	"fmt"
	"net"
	"runtime"
)

// Panic
func Panic(msg string, err error) {
	if err != nil {
		panic(errors.New(fmt.Sprintf("%s Error: %v\n", msg, err)))
	}
}

// MergeMaps merge items
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

		// Check the IP address to determine whether it is a loopback address
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

