package curl

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var address = "http://localhost:8080/check"

func TestGet(t *testing.T) {
	response, err := Get(address)
	assert.Nil(t, err)
	fmt.Println(response.String())
}

func TestPost(t *testing.T) {
	response, err := Post(address, nil)
	assert.Nil(t, err)
	fmt.Println(response.String())
}

func TestPut(t *testing.T) {
	response, err := Put(address, nil)
	assert.Nil(t, err)
	fmt.Println(response.String())
}

func TestPatch(t *testing.T) {
	response, err := Patch(address, nil)
	assert.Nil(t, err)
	fmt.Println(response.String())
}

func TestDelete(t *testing.T) {
	response, err := Delete(address)
	assert.Nil(t, err)
	fmt.Println(response.String())
}

func TestResponse(t *testing.T) {

	request := NewRequest("xxx", address, nil)
	if err := request.Err(); err != nil {
		t.Fatal(request.Err())

	}
	request.Header.Set("token", "123")
	response, err := request.Send()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(response.String())
	fmt.Println(string(response.Byte()))

	respObj := struct {
		Code    int    `json:"status_code"`
		Message string `json:"message"`
	}{}
	assert.Nil(t, response.BindJson(&respObj))
	fmt.Println(respObj.Message)

}

func TestPostFile(t *testing.T) {
	params := map[string]string{"name": "aa"}
	files := map[string]string{"file": "/usr/local/aa.file"}

	resp, err := PostFile(address, params, files)
	assert.Nil(t, err)
	fmt.Print(resp.String())
}
