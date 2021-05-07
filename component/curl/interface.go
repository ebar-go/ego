package curl

import (
	"io"
	"net/http"
)

// Curl send http request with sample func
type Curl interface {
	// Get send get request
	Get(url string) (Response, error)
	// Post send post request
	Post(url string, body io.Reader) (Response, error)
	// Put send put request
	Put(url string, body io.Reader) (Response, error)
	// Patch send patch request
	Patch(url string, body io.Reader) (Response, error)
	// Delete send delete request
	Delete(url string) (Response, error)
	// PostFile send post request with file
	PostFile(url string, files map[string]string, params map[string]string) (Response, error)
	// Send send origin request
	Send(request *http.Request) (Response, error)
}
