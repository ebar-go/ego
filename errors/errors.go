package errors

import (
	"fmt"
	"github.com/ebar-go/ego/utils/json"
	"net/http"
)

// Error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	MysqlConnectFailedCode = 1001
	RedisConnectFailedCode = 1002
)

// Error string
func (e *Error) Error() string {
	s, _ := json.Encode(e)
	return s
}

// New
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Format 格式化输出
func Format(code int, format string, v ...interface{}) *Error {
	return New(code, fmt.Sprintf(format, v...))
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(errStr string) *Error {
	e := new(Error)

	if err := json.Decode([]byte(errStr), e); err != nil {
		e.Code = http.StatusInternalServerError
		e.Message = err.Error()
	}
	return e
}

// Unauthorized generates a 401 error.
func Unauthorized(format string, v ...interface{}) *Error {
	return New(http.StatusUnauthorized, fmt.Sprintf(format, v...))
}

// Forbidden generates a 403 error.
func Forbidden(format string, v ...interface{}) *Error {
	return New(http.StatusForbidden, fmt.Sprintf(format, v...))
}

// NotFound generates a 404 error.
func NotFound(format string, v ...interface{}) *Error {
	return New(http.StatusNotFound, fmt.Sprintf(format, v...))
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(format string, v ...interface{}) *Error {
	return New(http.StatusMethodNotAllowed, fmt.Sprintf(format, v...))
}

// Timeout generates a 408 error.
func Timeout(format string, v ...interface{}) *Error {
	return New(http.StatusRequestTimeout, fmt.Sprintf(format, v...))
}

// InternalServerError generates a 500 error.
func InternalServer(format string, v ...interface{}) *Error {
	return New(http.StatusInternalServerError, fmt.Sprintf(format, v...))
}
