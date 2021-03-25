package curl

import "time"

type options struct {
	timeout time.Duration
}

type Option interface {
	apply(options *options)
}

type timeoutOption time.Duration

func (o timeoutOption) apply(options *options) {
	options.timeout = time.Duration(o)
}

func WithTimeout(timeout time.Duration) Option {
	return timeoutOption(timeout)
}