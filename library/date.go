package library

import (
	"time"
	"fmt"
)

const(
	defaultTimeFormat = "2006-01-02 15:04:05"
	defaultDateFormat = "2006-01-02"
)

// GetDateStr 获取日期字符串
func GetDateStr() string {
	return time.Now().Local().Format(defaultDateFormat)
}

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

// GetLastMonthTimeStr 获取上个月的时间
func GetLastMonthTimeStr() string {
	return time.Now().Local().AddDate(0, -1, 0).Format(defaultTimeFormat)
}

func GetTimeStampFloatStr() string {
	return fmt.Sprintf("%.6f", float64(time.Now().Local().UnixNano()) / 1e9)
}
