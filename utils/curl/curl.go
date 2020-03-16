package curl

import (
	"errors"
	"github.com/ebar-go/ego/app"
	"io/ioutil"
	"net/http"
)


// Execute return response as string
func Execute(req *http.Request) (string, error) {
	resp, err := app.Http().Do(req)
	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", errors.New("no response")
	}

	return stringify(resp)
}

// stringify stringify http response, return string and error
func stringify(resp *http.Response) (string, error) {
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("response status code not 200")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// close response
	defer func() {
		_ = resp.Body.Close()
	}()

	return string(data), nil
}
