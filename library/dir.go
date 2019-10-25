package library

import "os"

// GetCurrentPath 获取当前路径
func GetCurrentPath() string {
	path, _ := os.Getwd()
	return path
}
