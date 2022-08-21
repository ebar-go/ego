package component

import (
	"bytes"
	"errors"
	"fmt"
	serializer2 "github.com/ebar-go/ego/utils/serializer"
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

// Curl simple wrapper of http.Client
type Curl struct {
	Named
	httpClient   *http.Client
	bufferLenght int
}

// Get send get request
func (c *Curl) Get(url string) (serializer2.Serializer, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Send(request)
}

// Post send post request
func (c *Curl) Post(url string, body io.Reader) (serializer2.Serializer, error) {
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return c.Send(request)
}

// Put send put request
func (c *Curl) Put(url string, body io.Reader) (serializer2.Serializer, error) {
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	return c.Send(request)
}

// Delete send delete request
func (c *Curl) Delete(url string, body io.Reader) (serializer2.Serializer, error) {
	request, err := http.NewRequest(http.MethodDelete, url, body)
	if err != nil {
		return nil, err
	}
	return c.Send(request)
}

// PostFile send post request with file
func (c *Curl) PostFile(url string, files map[string]string, params map[string]string) (serializer2.Serializer, error) {
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

	return c.Send(request)
}

// Send return Response by http.Request
func (c *Curl) Send(request *http.Request) (serializer2.Serializer, error) {
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

func (c *Curl) ReadResponse(resp *http.Response) (serializer2.Serializer, error) {
	if resp == nil {
		return nil, EmptyResponse
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is:%d", resp.StatusCode)
	}

	length := int(resp.ContentLength)
	if length == 0 {
		length = c.bufferLenght
	}
	buffer := serializer2.NewBuffer(length)

	_, err := io.Copy(buffer, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return buffer, nil
}

func NewCurl() *Curl {
	return &Curl{
		Named:        Named("curl"),
		bufferLenght: 512,
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
