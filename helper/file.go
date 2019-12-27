package helper

import (
	"os"
)

//IsPathExist 判断文件或文件夹是否存在
func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}
	return true
}
