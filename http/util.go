package http

import (
	"net/http"
	"io/ioutil"
	"github.com/pkg/errors"
)

// 将response序列化
func  StringifyResponse(response *http.Response) (string, error) {
	if response == nil {
		return "", errors.New("没有响应数据")
	}

	if response.StatusCode != 200 {
		return "", errors.New("非200的上游返回")
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.WithMessage(err, "读取响应数据失败:")
	}

	// 关闭响应
	defer func() {
		response.Body.Close()
	}()

	return string(data), nil
}
