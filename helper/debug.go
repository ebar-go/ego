package helper

import (
	"runtime"
	"fmt"
)

// Debug 打印信息
func Debug(params ...interface{})  {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("[Trace]%s[%d]:%v \n", file, line, params)
	}
}

// Trace 返回trace日志
func Trace() []string {
	trace := []string{}
	for i:=0; i<10; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			trace = append(trace, fmt.Sprintf("[Trace]%s[%d]: \n", file, line))
		}
	}

	return trace
}
