package library

import "time"

const(
	defaultTimeFormat = "2006-01-02 15:04:05"
)

// GetTimeStr 获取时间字符串
func GetTimeStr() string {
	return time.Now().Local().Format(defaultTimeFormat)
}

// GetDefaultTimeFormat 获取默认时间格式
func GetDefaultTimeFormat() string {
	return defaultTimeFormat
}

// GetTimeStamp 获取时间戳
func GetTimeStamp() int64 {
	return time.Now().Local().Unix()
}


