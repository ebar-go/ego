package file

import (
	"os"
	"path/filepath"
)

// GetCurrentPath return compiled executable file absolute path
func GetCurrentPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// GetCurrentDir return compiled executable file directory
func GetCurrentDir() string {
	return filepath.Dir(GetCurrentPath())
}

//Exist check the given path exists
func Exist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}