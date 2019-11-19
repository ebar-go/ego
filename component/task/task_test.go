package task

import (
	"testing"
	"fmt"
)

// 测试初始化任务管理器
func TestInitManager(t *testing.T) {
	cron := InitManager()
	cron.AddFunc("*/5 * * * * ?", func() {
		fmt.Println("cron running")
	})

	cron.Start()
}
