/**
 * @Author: Hongker
 * @Description:
 * @File:  config
 * @Version: 1.0.0
 * @Date: 2021/4/3 18:13
 */

package http

type Config struct {
	Port int
	// jwt的key
	JwtSignKey []byte
	// trace header key
	TraceHeader string
	// 是否开启pprof
	Pprof bool
	// 是否开启swagger文档
	Swagger bool
}
