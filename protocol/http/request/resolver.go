package request

import (
	"github.com/ebar-go/ego/utils/serializer"
	"net/http"
)

// RequestInfoResolver represents a request info resolver that can be used to resolver requests.
type RequestInfoResolver interface {
	NewRequestInfo(req *http.Request) (*RequestInfo, error)
}

// RequestInfo includes information about a request.
type RequestInfo struct {
	serializer serializer.Serializer
}
