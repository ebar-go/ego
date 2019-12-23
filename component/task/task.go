// task 定时任务包
package task

import "github.com/robfig/cron"

// InitManager 初始化任务管理器
func InitManager() *cron.Cron {
	return cron.New()
}
