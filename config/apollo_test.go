package config

import (
	"testing"
	"fmt"
	"github.com/robfig/cron"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

func getApollo() *Apollo {
	return &Apollo{
		AppId: "open-api",
		Cluster: "local",
		Ip: "192.168.0.19:8080",
		Namespace: "application",
	}
}

// TestInitApolloConfig 测试初始化
func TestApollo_Init(t *testing.T) {
	apollo := getApollo()
	err := apollo.Init()
	assert.Nil(t, err)


}

// TestApollo_ListenChangeEvent 测试监听配置变更
func TestApollo_ListenChangeEvent(t *testing.T) {
	apollo := getApollo()
	c := cron.New()
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		changeEvent := apollo.ListenChangeEvent()

		bytes, _ := json.Marshal(changeEvent)

		fmt.Println("cron running:")
		fmt.Println("event:", string(bytes))
	})
	c.Start()

}

// TestApollo_GetStringValue 测试获取配置
func TestApollo_GetStringValue(t *testing.T) {
	apollo := getApollo()
	err := apollo.Init()
	assert.Nil(t, err)

	fmt.Println(apollo.GetStringValue("LOG_FILE",""))
}
