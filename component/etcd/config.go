/**
 * @Author: Hongker
 * @Description:
 * @File:  conf
 * @Version: 1.0.0
 * @Date: 2020/6/17 21:26
 */

package etcd

type Config struct {
	Endpoints []string
	Timeout   int
}
