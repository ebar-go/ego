package constant

const (
	TraceID        = "trace_id"      // 全局trace_id
	GatewayTrace   = "gateway-trace" // 网关trace
	JwtTokenMethod = "Bearer"
	JwtTokenHeader = "Authorization"
	JwtUserKey     = "jwt_user"
	JwtExpiredTime = 600

)

const (
	DefaultLogPath = "/wwwlogs/"
	DefaultSystemName = "system_name"
	SystemLogPrefix = "system_"
	RequestLogPrefix = "request_"
	AppLogPrefix = "app_"
	LogSuffix = ".log"
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
