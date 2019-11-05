package request

import (
	"testing"
	"github.com/ebar-go/ego/library"
)

func TestKongRequest_NewRequest(t *testing.T) {
	kong := Kong{
		Iss:"common-openapi",
		Secret:"Bf8DetEqOw4DOePtdOISWrwIyyboKH7h",
		Address: "http://app-gateway.internal.epetbar.com:8000",
	}

	request, _ := kong.NewRequest("gott-wms","GET", "/v1/basicInformation/warehouse/list?ware_nos=163", nil)

	resp, err := DefaultClient().Do(request)
	library.Debug(resp, err)
	str, _ :=StringifyResponse(resp)
	library.Debug(str)
}
