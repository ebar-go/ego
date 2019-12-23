package config

import (
	"github.com/ebar-go/ego/helper"
)


// Options 系统配置项
type Config struct {
	// 服务名称
	ServiceName string

	// 服务端口号
	ServicePort int

	// 响应日志最大长度
	MaxResponseLogSize int

	// 日志路径
	LogPath string

	// jwt的key
	JwtSignKey []byte
}

// init 通过读取环境变量初始化系统配置
func NewInstance() *Config {
	instance := &Config{}
	instance.ServiceName = helper.DefaultString(Getenv("SERVICE_NAME"), "app")
	instance.ServicePort = helper.DefaultInt(helper.String2Int(Getenv("SERVICE_PORT")), 8080)

	instance.LogPath = helper.DefaultString(Getenv("LOG_PATH"), "/tmp")
	instance.MaxResponseLogSize = helper.DefaultInt(helper.String2Int(Getenv("MAX_RESPONSE_LOG_SIZE")), 1000)

	instance.JwtSignKey = []byte(Getenv("JWT_SIGN_KEY"))

	return instance
}
