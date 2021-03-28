package curl

import (
	"bytes"
	"fmt"
	"github.com/ebar-go/egu"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)


type Curl interface {
	Get(url string) (Response, error)
	Post(url string, body io.Reader) (Response, error)
	Put(url string, body io.Reader) (Response, error)
	Patch(url string, body io.Reader) (Response, error)
	Delete(url string) (Response, error)
	PostFile(url string, files map[string]string, params map[string]string) (Response, error)
	Send(request *http.Request) (Response, error)
}

type curl struct {
	httpClient *http.Client
	pool *egu.BufferPool
}

func New(opts ...Option) Curl {
	options := options{timeout: time.Second * 3}
	for _, option := range opts {
		option.apply(&options)
	}
	return &curl{
		httpClient: &http.Client{
			Transport: &http.Transport{ // 配置连接池
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				IdleConnTimeout: 3 * time.Second,
			},
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       options.timeout,
		},
		pool:       egu.NewBufferPool(),
	}
}

type Response interface {
	String() string
	Byte() []byte
	BindJson(object interface{}) error
	Reader() io.Reader
}


// Get
func (c *curl) Get(url string) (Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Send(request)
}

// Post
func (c *curl) Post(url string, body io.Reader) (Response, error) {
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return c.Send(request)
}

// Put
func (c *curl) Put(url string, body io.Reader) (Response, error) {
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	return c.Send(request)
}

// Patch
func (c *curl) Patch(url string, body io.Reader) (Response, error) {
	request, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}
	return c.Send(request)
}

// Delete
func (c *curl) Delete(url string) (Response, error) {
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Send(request)
}

// PostFile 上传文件
func (c *curl) PostFile(url string, files map[string]string, params map[string]string) (Response, error) {
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

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	return c.Send(request)
}

func (c *curl) Send(request *http.Request) (Response, error) {
	resp, err := c.httpClient.Do(request)
	// close response
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}

	}()

	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("no response")
	}



	body, err := c.pool.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	return &response{body: body}, nil
}
