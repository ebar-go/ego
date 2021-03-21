/**
 * @Author: Hongker
 * @Description:
 * @File:  options
 * @Version: 1.0.0
 * @Date: 2021/3/21 21:13
 */

package http

import "github.com/gin-gonic/gin"

type options struct {
	// http port
	port int
	notFoundHandler gin.HandlerFunc
}

type Option interface {
	apply(opts *options)
}

type portOption int

func (o portOption) apply(opts *options) {
	opts.port = int(o)
}
func WithPort(port int) Option  {
	return portOption(port)
}
type notFoundOption gin.HandlerFunc
func (o notFoundOption) apply(opts *options) {
	opts.notFoundHandler = gin.HandlerFunc(o)
}
func With404Handler(handler gin.HandlerFunc) Option  {
	return notFoundOption(handler)
}
