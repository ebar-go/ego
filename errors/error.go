package errors

import "fmt"

const (
	unknown        = 200
	invalidParam   = 400
	unauthorized   = 401
	forbidden      = 403
	notfound       = 404
	timeout        = 405
	internalServer = 500
)

type Error struct {
	code    int
	message string
}

func (err Error) Error() string {
	return fmt.Sprintf("code=%d message=%s", err.code, err.message)
}

func (err Error) Code() int {
	return err.code
}

func (err Error) Message() string {
	return err.message
}

func (err *Error) WithMessage(message string) *Error {
	if err == nil {
		return nil
	}
	err.message = fmt.Sprintf("%s: %s", message, err.message)
	return err
}

func WithMessage(err error, message string) *Error {
	e := Convert(err)
	return e.WithMessage(message)
}

func Convert(err error) *Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		return e
	}

	return &Error{message: err.Error()}

}

func Sprintf(code int, format string, args ...interface{}) *Error {
	return &Error{code: code, message: fmt.Sprintf(format, args...)}
}
func Unknown(format string, args ...interface{}) *Error {
	return Sprintf(unknown, format, args...)
}

func NotFound(format string, args ...interface{}) *Error {
	return Sprintf(notfound, format, args...)
}

func InvalidParam(format string, args ...interface{}) *Error {
	return Sprintf(invalidParam, format, args...)
}
func Timeout(format string, args ...interface{}) *Error {
	return Sprintf(timeout, format, args...)
}
func Forbidden(format string, args ...interface{}) *Error {
	return Sprintf(forbidden, format, args...)
}
func Unauthorized(format string, args ...interface{}) *Error {
	return Sprintf(unauthorized, format, args...)
}
func InternalServer(format string, args ...interface{}) *Error {
	return Sprintf(internalServer, format, args...)
}
