package buffer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// Pool buffer instance
type Pool struct {
	instance sync.Pool
}

// NewPool
func NewPool() *Pool {
	return &Pool{instance: sync.Pool{New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 4096))
	}}}
}

// ReadResponse read http response
func (p *Pool) ReadResponse(response *http.Response) ([]byte, error) {
	if response == nil {
		return nil, fmt.Errorf("response is empty")
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is:%d", response.StatusCode)
	}

	buffer := p.instance.Get().(*bytes.Buffer)
	buffer.Reset()
	defer func() {
		if buffer != nil {
			p.instance.Put(buffer)
			buffer = nil
		}
	}()
	_, err := io.Copy(buffer, response.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read respone:%s", err.Error())
	}

	return buffer.Bytes(), nil
}
