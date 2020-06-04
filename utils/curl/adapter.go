package curl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// Adapter buffer pool
type Adapter struct {
	pool sync.Pool
}

// NewAdapter
func NewAdapter() *Adapter {
	return &Adapter{pool: sync.Pool{New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 4096))
	}}}
}

// StringifyResponse return response body as string
func (adapter *Adapter) StringifyResponse(response *http.Response) (string, error) {
	if response == nil {
		return "", fmt.Errorf("response is empty")
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response status code is:%d", response.StatusCode)
	}

	buffer := adapter.pool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer func() {
		if buffer != nil {
			adapter.pool.Put(buffer)
			buffer = nil
		}
	}()
	_, err := io.Copy(buffer, response.Body)

	if err != nil {
		return "", fmt.Errorf("failed to read respone:%s", err.Error())
	}

	return buffer.String(), nil
}
