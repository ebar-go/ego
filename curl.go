package ego

import (
	"fmt"
	"github.com/ebar-go/ego/app"
	"net/http"
)

// Curl send http request
func Curl(req *http.Request) (string, error) {
	resp, err := app.Http().Do(req)
	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", fmt.Errorf("no response")
	}

	// close response
	defer func() {
		_ = resp.Body.Close()
	}()

	return app.BufferPool().StringifyResponse(resp)
}
