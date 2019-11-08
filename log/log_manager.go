package log

import (
	"os"
	"fmt"
	"github.com/ebar-go/ego/library"
	"path/filepath"
)

var AppLogger *Logger
var SystemLogger *Logger
var RequestLogger *Logger



// App 应用日志
func App() *Logger {
	if AppLogger == nil {
		AppLogger = New()
	}

	return AppLogger
}

// System 系统日志
func System() *Logger {
	if SystemLogger == nil {
		SystemLogger = New()
	}

	return SystemLogger
}

// Request 请求日志
func Request() *Logger {
	if RequestLogger == nil {
		RequestLogger = New()
	}

	return RequestLogger
}

// NewFileLogger 新的文件日志管理器
func NewFileLogger(filePath string) *Logger {
	logger := New()

	if !library.IsPathExist(filePath) {
		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil{
			library.Debug(err)
			return logger
		}
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err == nil {
		logger.SetOutWriter(file)
		fmt.Printf("初始化日志文件:%s,成功\n", filePath)
	}else {
		fmt.Println("err:" + err.Error())
	}

	return logger
}
