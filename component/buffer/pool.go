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

// StringifyResponse return response body as string
func (p *Pool) StringifyResponse(response *http.Response) (string, error) {
	if response == nil {
		return "", fmt.Errorf("response is empty")
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response status code is:%d", response.StatusCode)
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
		return "", fmt.Errorf("failed to read respone:%s", err.Error())
	}

	return buffer.String(), nil
}
