/**
 * @Author: Hongker
 * @Description:
 * @File:  conf
 * @Version: 1.0.0
 * @Date: 2020/6/17 21:26
 */

package etcd

import "time"

// Config 配置项
type Config struct {
	// 节点
	Endpoints []string
	// 超时时间
	Timeout   time.Duration
}
