/**
 * @Author: Hongker
 * @Description:
 * @File:  config
 * @Version: 1.0.0
 * @Date: 2021/4/3 18:15
 */

package log

type Config struct {
	// 日志路径
	Path string

	// 是否开启debug,开启后会显示debug信息
	Debug bool
}
