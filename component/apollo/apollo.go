/**
Apollo配置初始化、监听配置变动、获取配置
*/
package apollo

import (
	"github.com/zouyx/agollo"
	"os"
)

// Conf apollo config
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

// Init
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

	if err := agollo.Start(); err != nil {
		return err
	}

	loadEnvironment()
	return nil
}

// loadEnvironment
func loadEnvironment()  {
	cache := agollo.GetApolloConfigCache().NewIterator()
	for {
		item := cache.Next()
		if item == nil {
			break
		}
		_ = os.Setenv(string(item.Key), string(item.Value))
	}
}

// ListenApolloChangeEvent listen apollo config change event
func ListenChangeEvent() <-chan *agollo.ChangeEvent {
	return agollo.ListenChangeEvent()
}

// GetStringValue get config as string
func GetStringValue(key, defaultValue string) string {
	return agollo.GetStringValue(key, defaultValue)
}

// GetIntValue get config as int
func GetIntValue(key string, defaultValue int) int {
	return agollo.GetIntValue(key, defaultValue)
}

// GetBoolValue get config as bool
func GetBoolValue(key string, defaultValue bool) bool {
	return agollo.GetBoolValue(key, defaultValue)
}

// GetFloatValue get config as float
func GetFloatValue(key string, defaultV float64) float64 {
	return agollo.GetFloatValue(key, defaultV)
}
