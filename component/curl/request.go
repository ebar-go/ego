package curl

import (
	"bytes"
	"fmt"
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/egu"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var bufferPool = egu.NewBufferPool()

type request struct {
	*http.Request
	err error
}

// NewRequest
func NewRequest(method, url string, body io.Reader) *request {
	req := new(request)
	request, err := http.NewRequest(method, url, body)

	if err != nil {
		req.err = err
	}

	req.Request = request
	return req
}

// Get
func Get(url string) (*response, error) {
	return NewRequest(http.MethodGet, url, nil).Send()
}

// Post
func Post(url string, body io.Reader) (*response, error) {
	return NewRequest(http.MethodPost, url, body).Send()
}

// Put
func Put(url string, body io.Reader) (*response, error) {
	return NewRequest(http.MethodPut, url, body).Send()
}

// Patch
func Patch(url string, body io.Reader) (*response, error) {
	return NewRequest(http.MethodPatch, url, body).Send()
}

// Delete
func Delete(url string) (*response, error) {
	return NewRequest(http.MethodDelete, url, nil).Send()
}

// PostFile 上传文件
func PostFile(url string, files map[string]string, params map[string]string) (*response, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	// 添加form参数
	for name, value := range params {
		_ = writer.WriteField(name, value)
	}

	// 写入文件流
	for field, path := range files {
		// 读取文件
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("Open File: %v", err)
		}
		_ = file.Close()

		// 写入writer
		part, err := writer.CreateFormFile(field,filepath.Base(path))
		if err != nil {
			return nil, fmt.Errorf("Create Form File: %v", err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, err
		}
	}

	// 必须close，这样writer.FormDataContentType()才正确
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("Close Writer: %v", err)
	}

	req := NewRequest(http.MethodPost, url, body)
	if req.err != nil {
		return nil, req.err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req.Send()
}

// Send send http request
func (req *request) Send() (*response, error) {
	if req.err != nil {
		return nil, req.err
	}
	req.Header.Set(app.Config().Server().TraceHeader, trace.Get())
	resp, err := app.Http().Do(req.Request)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("no response")
	}

	// close response
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := bufferPool.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	return &response{body: body}, nil
}

func (req *request) Err() error {
	return req.err
}
