package curl

import "time"

// Option curl option interface
type Option interface {
	apply(options *options)
}

// options 配置项
type options struct {
	// 超时
	timeout time.Duration
}

// timeoutOption set http client timeout option
type timeoutOption time.Duration

func (o timeoutOption) apply(options *options) {
	options.timeout = time.Duration(o)
}
// WithTimeout use timeout option
func WithTimeout(timeout time.Duration) Option {
	return timeoutOption(timeout)
}
