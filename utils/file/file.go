package file

import (
	"os"
	"path/filepath"
	"runtime"
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

// GetExecuteDir return current executable file directory,different with GetCurrentDir
func GetExecuteDir() string {
	_, fileStr, _, _ := runtime.Caller(0)
	return filepath.Dir(fileStr)
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

// Mkdir Create the DIRECTORY(ies), if they do not already exist
// parents no error if existing, make parent directories as needed
func Mkdir(dir string, parents bool) error {
	if Exist(dir) {
		return nil
	}

	if parents {
		return os.MkdirAll(dir, os.ModePerm)
	}

	return os.Mkdir(dir, os.ModePerm)


}