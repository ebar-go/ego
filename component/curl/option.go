package curl

import "time"

type Option interface {
	apply(options *options)
}

// options 配置项
type options struct {
	// 超时
	timeout time.Duration
}

type timeoutOption time.Duration

func (o timeoutOption) apply(options *options) {
	options.timeout = time.Duration(o)
}

func WithTimeout(timeout time.Duration) Option {
	return timeoutOption(timeout)
}
