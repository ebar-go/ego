package curl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// ResponseReader read http response body
type ResponseReader interface {
	Read(response *http.Response) ([]byte, error)
}

// bufferPoolResponseReader
type bufferPoolResponseReader struct {
	pool sync.Pool
}

// NewResponseReader
func NewResponseReader() *bufferPoolResponseReader {
	return &bufferPoolResponseReader{pool: sync.Pool{New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 4096))
	}}}
}

// Read read http response
func (p *bufferPoolResponseReader) Read(response *http.Response) ([]byte, error) {
	if response == nil {
		return nil, fmt.Errorf("response is empty")
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is:%d", response.StatusCode)
	}

	buffer := p.pool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer func() {
		if buffer != nil {
			p.pool.Put(buffer)
			buffer = nil
		}
	}()
	_, err := io.Copy(buffer, response.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read respone:%s", err.Error())
	}

	return buffer.Bytes(), nil
}
