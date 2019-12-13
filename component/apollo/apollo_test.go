package apollo

import (
	"testing"
	"fmt"
	"github.com/robfig/cron"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/ebar-go/ego/component/config"
)

// prepareConf
func prepareConf() Conf {
	var conf = struct {
		Apollo struct{
			AppId string `yaml:"appId"`
			Cluster string
			Namespace string
			Ip string
			BackupConfigPath string `yaml:"backConfigPath"`
		}
	}{}


	if err := config.LoadYaml(&conf, "/tmp/app.yaml"); err != nil {
		panic(err)
	}

	return Conf{
		AppId: conf.Apollo.AppId,
		Cluster: conf.Apollo.Cluster,
		Ip: conf.Apollo.Ip,
		Namespace: conf.Apollo.Namespace,
		BackupConfigPath: conf.Apollo.BackupConfigPath,
	}
}

// TestInit 测试初始化
func TestInit(t *testing.T) {
	err := Init(prepareConf())
	assert.Nil(t, err)
}

// TestListenChangeEvent 测试监听配置变更
func TestListenChangeEvent(t *testing.T) {
	c := cron.New()
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		changeEvent := ListenChangeEvent()

		bytes, _ := json.Marshal(changeEvent)

		fmt.Println("cron running:")
		fmt.Println("event:", string(bytes))
	})
	c.Start()
}

// TestGetStringValue 测试获取配置
func TestGetStringValue(t *testing.T) {

	fmt.Println(GetStringValue("LOG_FILE",""))
}

// TestGetIntValue
func TestGetIntValue(t *testing.T) {
	fmt.Println(GetIntValue("HTTP_PORT",8080))
}

// TestGetBoolValue
func TestGetBoolValue(t *testing.T) {
	fmt.Println(GetBoolValue("APP_DEBUG", false))
}

// TestMain main
func TestMain(m *testing.M) {
	fmt.Println(prepareConf())
	Init(prepareConf())
	m.Run()
}
