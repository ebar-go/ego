package response

import "sync"

// 默认对象池
var rp = newPool()

// responsePool 响应内容对象池
type responsePool struct {
	pool sync.Pool
}

// newPool
func newPool() *responsePool {
	return &responsePool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(response)
			},
		},
	}
}

// Get
func (p *responsePool) Get() *response {
	r := p.pool.Get().(*response)
	r.Reset()
	defer func() {
		if r != nil {
			p.pool.Put(r)
			r = nil
		}
	}()
	return r
}
