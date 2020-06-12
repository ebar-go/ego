package ego

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestExecute(t *testing.T) {
	address := "http://baidu.com"

	request, err := http.NewRequest(http.MethodGet, address, nil)
	assert.Nil(t, err)

	response, err := Curl(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
	fmt.Println(response)
}
