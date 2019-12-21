package config

import (
	"github.com/ebar-go/ego/helper"
)

// 系统日志实例
var Instance = new(SystemConfig)

// SystemConfig 系统日志
type SystemConfig struct {
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
func init()  {
	Instance.ServiceName = helper.DefaultString(Getenv("SERVICE_NAME"), "app")
	Instance.ServicePort = helper.DefaultInt(helper.String2Int(Getenv("SERVICE_PORT")), 8080)

	Instance.LogPath = helper.DefaultString(Getenv("LOG_PATH"), "/tmp")
	Instance.MaxResponseLogSize = helper.DefaultInt(helper.String2Int(Getenv("MAX_RESPONSE_LOG_SIZE")), 1000)

	Instance.JwtSignKey = []byte(Getenv("JWT_SIGN_KEY"))
}
