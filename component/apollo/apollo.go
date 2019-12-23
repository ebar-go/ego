/**
Apollo配置初始化、监听配置变动、获取配置
*/
package apollo

import (
	"github.com/zouyx/agollo"
)

// Apollo apollo配置项
type Conf struct {
	// application id
	AppId string `json:"appId"`

	// apollo cluster
	Cluster string `json:"cluster"`

	// apollo application
	Namespace string `json:"namespaceName"`

	// server address ex:127.0.0.1:8080
	Ip string `json:"ip"`

	// apolloConfig.json path
	BackupConfigPath string `json:"backup_config_path"`
}

// Init 初始化apollo配置
func Init(conf Conf) error {
	agollo.InitCustomConfig(func() (*agollo.AppConfig, error) {
		return &agollo.AppConfig{
			AppId:            conf.AppId,
			Cluster:          conf.Cluster,
			Ip:               conf.Ip,
			NamespaceName:    conf.Namespace,
			BackupConfigPath: conf.BackupConfigPath,
		}, nil
	})

	return agollo.Start()
}

// ListenApolloChangeEvent 监听配置变动
func ListenChangeEvent() <-chan *agollo.ChangeEvent {
	return agollo.ListenChangeEvent()
}

// GetStringValue 获取字符串配置
func GetStringValue(key, defaultValue string) string {
	return agollo.GetStringValue(key, defaultValue)
}

// GetIntValue 获取整型配置
func GetIntValue(key string, defaultValue int) int {
	return agollo.GetIntValue(key, defaultValue)
}

// GetBoolValue 获取bool配置
func GetBoolValue(key string, defaultValue bool) bool {
	return agollo.GetBoolValue(key, defaultValue)
}

// GetFloatValue 获取浮点型配置
func GetFloatValue(key string, defaultV float64) float64 {
	return agollo.GetFloatValue(key, defaultV)
}
