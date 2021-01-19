package server


// Conf  服务配置
type Config struct {
	// 运行环境
	Environment string
	// 服务名称
	Name string
	// 服务端口号,
	Port int
	// 响应日志最大长度
	MaxResponseLogSize int
	// 日志路径
	LogPath string
	// jwt的key
	JwtSignKey []byte
	// trace header key
	TraceHeader string
	// http request timeout
	HttpRequestTimeOut int
	// 是否开启debug,开启后会显示debug信息
	Debug bool
	// 是否开启pprof
	Pprof bool
	// 是否开启swagger文档
	Swagger bool
	// 是否开启定时任务
	Task bool
}
