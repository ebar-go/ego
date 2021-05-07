package curl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// httpClientImpl instance of Curl interface
type httpClientImpl struct {
	httpClient *http.Client
	reader       ResponseReader
}



// Get send get request
func (impl *httpClientImpl) Get(url string) (Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return impl.Send(request)
}

// Post send post request
func (impl *httpClientImpl) Post(url string, body io.Reader) (Response, error) {
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return impl.Send(request)
}

// Put send put request
func (impl *httpClientImpl) Put(url string, body io.Reader) (Response, error) {
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	return impl.Send(request)
}

// Patch send patch request
func (impl *httpClientImpl) Patch(url string, body io.Reader) (Response, error) {
	request, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}
	return impl.Send(request)
}

// Delete send delete request
func (impl *httpClientImpl) Delete(url string) (Response, error) {
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return impl.Send(request)
}

// PostFile send post request with file
func (impl *httpClientImpl) PostFile(url string, files map[string]string, params map[string]string) (Response, error) {
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

	return impl.Send(request)
}

// Send send http request 
func (impl *httpClientImpl) Send(request *http.Request) (Response, error) {
	resp, err := impl.httpClient.Do(request)
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
		return nil, errors.New("no response")
	}

	body, err := impl.reader.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	return &response{body: body}, nil
}


// New return a instance of Curl
func New(opts ...Option) Curl {
	options := options{timeout: time.Second * 30}
	for _, option := range opts {
		option.apply(&options)
	}
	return &httpClientImpl{
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
		reader: NewResponseReader(),
	}
}