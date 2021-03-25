package curl

import (
	"io"
	"net/http"
)


type Curl interface {
	Get(url string) (*response, error)
	Post(url string, body io.Reader) (*response, error)
	Put(url string, body io.Reader) (*response, error)
	Patch(url string, body io.Reader) (*response, error)
	Delete(url string) (*response, error)
	PostFile(url string, files map[string]string, params map[string]string) (*response, error)
	Send(request *http.Request) (*response, error)
}

type curl struct {


}
