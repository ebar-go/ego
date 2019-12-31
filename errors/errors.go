package errors

import (
	"fmt"
	"github.com/ebar-go/ego/utils/json"
)

// Error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error string
func (e *Error) Error() string {
	s, _ := json.Encode(e)
	return s
}

// New
func New( code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(errStr string) *Error {
	e := new(Error)

	if err := json.Decode([]byte(errStr), e); err != nil {
		e.Code = 500
		e.Message = err.Error()
	}
	return e
}

// Unauthorized generates a 401 error.
func Unauthorized(format string, v ...interface{}) *Error {
	return New(401, fmt.Sprintf(format, v...))
}

// Forbidden generates a 403 error.
func Forbidden(format string, v ...interface{}) *Error {
	return New(403, fmt.Sprintf(format, v...))
}

// NotFound generates a 404 error.
func NotFound(format string, v ...interface{}) *Error {
	return New(404, fmt.Sprintf(format, v...))
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(format string, v ...interface{}) *Error {
	return New(405, fmt.Sprintf(format, v...))
}

// Timeout generates a 408 error.
func Timeout(format string, v ...interface{}) *Error {
	return New(408, fmt.Sprintf(format, v...))
}

// InternalServerError generates a 500 error.
func InternalServerError(format string, v ...interface{}) *Error {
	return New(500, fmt.Sprintf(format, v...))
}
