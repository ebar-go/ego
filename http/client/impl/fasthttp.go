package impl

import (
	"github.com/valyala/fasthttp"
	"github.com/pkg/errors"
	"github.com/ebar-go/ego/http/client/object"
	"bytes"
)

type FastHttpClient struct {

}

func NewFastHttpClient() FastHttpClient {
	return FastHttpClient{}
}

func (client FastHttpClient) NewRequest(param object.RequestParam) object.IRequest{
	req := fasthttp.AcquireRequest()

	req.Header.SetContentType("application/json")
	req.Header.SetMethod(param.Method)

	req.SetRequestURI(param.Url)

	for key, value := range param.Headers {
		req.Header.Add(key, value)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(param.Body)
	if param.Body != nil {
		req.SetBody(buf.Bytes())
	}

	return req
}

// Send 发送http请求，得到响应
func (client FastHttpClient) Execute(request interface{}) (string, error) {
	if request == nil {
		return "", errors.New("request is nil")
	}

	req, ok := request.(*fasthttp.Request)
	if !ok {
		return "", errors.New("request not *http.request")
	}
	defer fasthttp.ReleaseRequest(req) // 用完需要释放资源

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源

	if err := fasthttp.Do(req, resp); err != nil {
		return "", errors.WithMessage(err, "请求失败")
	}

	return string(resp.Body()), nil
}
