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
	stack   *stack
}

func (err Error) Error() string {
	return fmt.Sprintf("code=%d message=%s stack=%v", err.code, err.message, err.stack.current())
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
	if err == nil {
		return nil
	}

	e, ok := err.(*Error)
	if !ok {
		e = sprintf(unknown, err.Error())
	}

	return e.WithMessage(message)
}

func Convert(err error) *Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		return e
	}

	return sprintf(unknown, err.Error())
}

func New(code int, msg string) *Error {
	return sprintf(code, msg)
}

func sprintf(code int, format string, args ...interface{}) *Error {
	return &Error{code: code, message: fmt.Sprintf(format, args...), stack: callers()}
}

func Unknown(format string, args ...interface{}) *Error {
	return sprintf(unknown, format, args...)
}

func NotFound(format string, args ...interface{}) *Error {
	return sprintf(notfound, format, args...)
}

func InvalidParam(format string, args ...interface{}) *Error {
	return sprintf(invalidParam, format, args...)
}
func Timeout(format string, args ...interface{}) *Error {
	return sprintf(timeout, format, args...)
}
func Forbidden(format string, args ...interface{}) *Error {
	return sprintf(forbidden, format, args...)
}
func Unauthorized(format string, args ...interface{}) *Error {
	return sprintf(unauthorized, format, args...)
}
func InternalServer(format string, args ...interface{}) *Error {
	return sprintf(internalServer, format, args...)
}

func Is(err error, target error) bool {
	e := Convert(err)
	if e == nil {
		return false
	}
	return e.code == Convert(target).code
}
