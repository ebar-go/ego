package errors

import (
	"fmt"
	"net/http"
)

// Error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error string
func (e *Error) Error() string {
	return e.Message
}

// New
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Sprintf 格式化输出
func Sprintf(code int, format string, v ...interface{}) *Error {
	return New(code, fmt.Sprintf(format, v...))
}

// With
func With(msg string, err error) *Error {
	if e, ok := err.(*Error); ok {
		e.Message = fmt.Sprintf("%s:%s", msg, e.Message)
		return e
	}
	return InternalServer(fmt.Sprintf("%s:%s", msg, err.Error()))
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
