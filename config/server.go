package config

// ServerOption 服务配置
type ServerOption struct {
	// 服务名称
	Name string

	// 服务端口号
	Port int

	// 响应日志最大长度
	MaxResponseLogSize int

	// 日志路径
	LogPath string
	// jwt的key
	JwtSignKey []byte

	// trace header key
	TraceHeader string
}
