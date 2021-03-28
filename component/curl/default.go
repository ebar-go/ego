package curl

import (
	"io"
	"net/http"
)


var _default = New()

func Default() Curl {
	return _default
}

func Get(url string) (Response, error) {
	return Default().Get(url)
}
func Post(url string, body io.Reader) (Response, error) {
	return Default().Post(url, body)
}
func Put(url string, body io.Reader) (Response, error) {
	return Default().Put(url, body)
}
func Patch(url string, body io.Reader) (Response, error) {
	return Default().Patch(url, body)
}
func Delete(url string) (Response, error) {
	return Default().Delete(url)
}
func PostFile(url string, files map[string]string, params map[string]string) (Response, error) {
	return Default().PostFile(url, files, params)
}

func Send(request *http.Request) (Response, error) {
	return Default().Send(request)
}