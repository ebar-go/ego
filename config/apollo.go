/**
Apollo配置初始化、监听配置变动、获取配置
 */
package config

import (
	"github.com/zouyx/agollo"
)

// Apollo apollo配置项
type Apollo struct {
	AppId string `json:"appId"`
	Cluster string `json:"cluster"`
	Namespace string `json:"namespaceName"`
	Ip string `json:"ip"`
}

// InitApollo 初始化apollo配置
func (apollo *Apollo) Init() error {
	agollo.InitCustomConfig(func () (*agollo.AppConfig, error) {
		return &agollo.AppConfig{
			AppId:         apollo.AppId,
			Cluster:       apollo.Cluster,
			Ip:            apollo.Ip,
			NamespaceName: apollo.Namespace,
		}, nil
	})

	return agollo.Start()
}

// ListenApolloChangeEvent 监听配置变动
func (apollo *Apollo) ListenChangeEvent() <-chan *agollo.ChangeEvent {
	return agollo.ListenChangeEvent()
}

// GetApolloStringValue 获取字符串配置
func (apollo *Apollo) GetStringValue(key , defaultValue string) string {
	return agollo.GetStringValue(key, defaultValue)
}