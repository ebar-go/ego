package request

import (
	"testing"
	"github.com/ebar-go/ego/library"
)

func TestKongRequest_NewRequest(t *testing.T) {
	kong := Kong{
		Iss:"xxx",
		Secret:"123",
		Address: "aaa",
	}

	request, _ := kong.NewRequest("user","GET", "test", nil)

	resp, err := DefaultClient().Do(request)
	library.Debug(resp, err)
	str, _ :=StringifyResponse(resp)
	library.Debug(str)
}
