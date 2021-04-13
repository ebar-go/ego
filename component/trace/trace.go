// Trace
// 基于goroutine的编号，实现共享全局唯一ID,用于串联业务线的上下文
package trace

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/petermattis/goid"
)

type Trace struct {
	// 使用并发map提高性能，原理是采用32个分片来分散
	ids cmap.ConcurrentMap
}

func NewTrace() *Trace {
	return &Trace{ids: cmap.New()}
}

// key 使用goroutine的id作为下标
func (t *Trace) key() string {
	return fmt.Sprintf("g%d", goid.Get())
}

// Set
func (t *Trace) Set(uuid string) {
	t.ids.Set(t.key(), uuid)
}

// Get
func (t *Trace) Get() string {
	id, exist := t.ids.Get(t.key())
	if exist {
		return id.(string)
	}
	return ""
}

// GC 回收
func (t *Trace) GC() {
	t.ids.Remove(t.key())
}

// Go 保持两个goroutine使用同一个traceId连通
func (t *Trace) Go(fn func()) {
	go func(traceId string) {
		t.Set(traceId)
		defer t.GC()
		fn()
	}(t.Get())
}
