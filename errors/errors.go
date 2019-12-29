package errors

import (
	"fmt"
	"github.com/ebar-go/ego/helper"
)

// Error
type Error struct {
	Code int `json:"code"`
	Key string `json:"key"`
	Message string `json:"message"`
}

// Error string
func (e *Error) Error() string {
	s, _ := helper.JsonEncode(e)
	return s
}

// New
func New(key, message string, code int) error {
	return &Error{
		Code:    code,
		Key:     key,
		Message: message,
	}
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(errStr string) *Error {
	e := new(Error)
	err := helper.JsonDecode([]byte(errStr), e)
	if err != nil {
		e.Code = 500
		e.Message = err.Error()
	}
	return e
}

// Unauthorized generates a 401 error.
func Unauthorized(key, format string, v ...interface{}) error {
	return New(key, fmt.Sprintf(format, v...), 401)
}

// Forbidden generates a 403 error.
func Forbidden(key, format string, v ...interface{}) error {
	return New(key, fmt.Sprintf(format, v...), 403)
}

// NotFound generates a 404 error.
func NotFound(key, format string, v ...interface{}) error {
	return New(key, fmt.Sprintf(format, v...), 404)
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(key, format string, v ...interface{}) error {
	return New(key, fmt.Sprintf(format, v...), 405)
}

// Timeout generates a 408 error.
func Timeout(key, format string, v ...interface{}) error {
	return New(key, fmt.Sprintf(format, v...), 408)
}

// InternalServerError generates a 500 error.
func InternalServerError(key, format string, v ...interface{}) error {
	return New(key, fmt.Sprintf(format, v...), 500)
}