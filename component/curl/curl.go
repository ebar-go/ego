package curl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ebar-go/ego/component/tracer"
	"github.com/ebar-go/ego/utils/serializer"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	EmptyResponse = errors.New("empty response")
)

// Instance simple wrapper of http.Client
type Instance struct {
	httpClient   *http.Client
	bufferLength int
	tracer       *tracer.Instance
}

// Get send get request
func (c *Instance) Get(ctx context.Context, url string) (serializer.Serializer, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Send(ctx, request)
}

// Post send post request
func (c *Instance) Post(ctx context.Context, url string, body io.Reader) (serializer.Serializer, error) {
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return c.Send(ctx, request)
}

// Put send put request
func (c *Instance) Put(ctx context.Context, url string, body io.Reader) (serializer.Serializer, error) {
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	return c.Send(ctx, request)
}

// Delete send delete request
func (c *Instance) Delete(ctx context.Context, url string, body io.Reader) (serializer.Serializer, error) {
	request, err := http.NewRequest(http.MethodDelete, url, body)
	if err != nil {
		return nil, err
	}
	return c.Send(ctx, request)
}

// PostFile send post request with file
func (c *Instance) PostFile(ctx context.Context, url string, files map[string]string, params map[string]string) (serializer.Serializer, error) {
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
			return nil, fmt.Errorf("open file: %v", err)
		}
		_ = file.Close()

		// 写入writer
		part, err := writer.CreateFormFile(field, filepath.Base(path))
		if err != nil {
			return nil, fmt.Errorf("create form file: %v", err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, err
		}
	}

	// 必须close，这样writer.FormDataContentType()才正确
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close writer: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	return c.Send(ctx, request)
}

// Send return Response by http.Request
func (c *Instance) Send(ctx context.Context, request *http.Request) (serializer.Serializer, error) {
	if c.tracer != nil {
		c.tracer.Set(c.tracer.Get())
		defer c.tracer.Release()
	}
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

	return c.ReadResponse(resp)
}

// ReadResponse reads body from the *http.Response
func (c *Instance) ReadResponse(resp *http.Response) (serializer.Serializer, error) {
	if resp == nil {
		return nil, EmptyResponse
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is:%d", resp.StatusCode)
	}

	length := int(resp.ContentLength)
	if length == 0 {
		length = c.bufferLength
	}
	buffer := serializer.NewBuffer(length)

	_, err := io.Copy(buffer, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return buffer, nil
}

// WithHttpClient sets the http client
func (c *Instance) WithHttpClient(httpClient *http.Client) *Instance {
	c.httpClient = httpClient
	return c
}

func (c *Instance) WithTracer(tracer *tracer.Instance) *Instance {
	c.tracer = tracer
	return c
}

func New() *Instance {
	return &Instance{
		bufferLength: 512,
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
			Timeout:       30 * time.Second,
		},
	}
}
