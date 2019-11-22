package request

import (
	"testing"
	"github.com/ebar-go/ego/helper"
)

func TestKongRequest_NewRequest(t *testing.T) {
	kong := Kong{
		Iss:"xxx",
		Secret:"123",
		Address: "aaa",
	}

	request, _ := kong.NewRequest("user","GET", "test", nil)

	resp, err := DefaultClient().Do(request)
	helper.Debug(resp, err)
	str, _ :=StringifyResponse(resp)
	helper.Debug(str)
}
