package library

import "time"

// GetTimeStr 获取时间字符串
func GetTimeStr() string {
	return time.Now().Local().Format("2006-01-02 15:04:05")
}