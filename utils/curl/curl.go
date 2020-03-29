package curl

import (
	"errors"
	"github.com/ebar-go/ego/app"
	"net/http"
)

// DefaultAdapter
var DefaultAdapter = NewAdapter()

// Execute return response as string
func Execute(req *http.Request) (string, error) {
	resp, err := app.Http().Do(req)
	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", errors.New("no response")
	}

	// close response
	defer func() {
		_ = resp.Body.Close()
	}()

	return DefaultAdapter.StringifyResponse(resp)
}