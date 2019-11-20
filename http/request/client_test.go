package request

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestDefaultClient(t *testing.T) {
	client := DefaultClient()
	assert.NotNil(t, client)
}

func TestNew(t *testing.T) {
	request, err := New(MethodGet, "https://www.baidu.com", nil)
	assert.Nil(t, err)
	assert.NotNil(t, request)

}

func TestClient_Do(t *testing.T) {
	client := DefaultClient()
	request, err := New(MethodGet, "https://www.baidu.com", nil)

	assert.Nil(t, err)
	resp, err := client.Do(request)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

}

func TestStringifyResponse(t *testing.T) {
	client := DefaultClient()
	request, err := New(MethodGet, "https://www.baidu.com", nil)

	assert.Nil(t, err)
	resp, err := client.Do(request)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	responseStr, err := StringifyResponse(resp)
	fmt.Println(responseStr)
}