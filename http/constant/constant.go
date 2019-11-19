package constant

const (
	TraceID        = "trace_id"      // 全局trace_id
	GatewayTrace   = "gateway-trace" // 网关trace
	JwtTokenMethod = "Bearer"
	JwtTokenHeader = "Authorization"
	JwtUserKey     = "jwt_user"
	JwtExpiredTime = 600
	SystemNameKey = "SYSTEM_NAME"

)

const (
	EnvSystemName = "SYSTEM_NAME"
	EnvSystemPort = "SYSTEM_PORT"
)

const (
	DefaultLogPath = "/wwwlogs/"
	DefaultSystemName = "system_name"
	SystemLogPrefix = "system-"
	RequestLogPrefix = "request-"
	AppLogPrefix = "app-"
	MqLogPrefix = "mq-"
	LogSuffix = ".log"

	AppLogComponentName = "app"
	TraceLogComponentName = "trace"
	MqLogComponentName = "mq"
	SystemLogComponentName = "phplogs"
)

const (
	TraceIdPrefix = "TraceId"
	RequestIdPrefix = "RequestId"
)

const (
	StatusOk           = 200
	StatusNotFound     = 404
	StatusUnauthorized = 401 // RFC 7235, 3.1
	StatusError        = 500
)

const (
	HttpSchema = "http://"
)
