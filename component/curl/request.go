package curl

import (
	"fmt"
	"github.com/zutim/ego/app"
	"github.com/zutim/ego/component/trace"
	"io"
	"net/http"
)

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

	bytes, err := app.BufferPool().ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	return &response{body: bytes}, nil
}

func (req *request) Err() error {
	return req.err
}
