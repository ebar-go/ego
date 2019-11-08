package helper

import (
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/ebar-go/ego/library"
	"github.com/ebar-go/ego/http/constant"
	"io/ioutil"
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
)


// GetTraceId 获取唯一traceId
func GetTraceId(c *gin.Context) string {
	traceIdInterface, exist := c.Get(constant.TraceID)
	traceId := ""
	if exist == false {
		traceId = c.GetHeader(constant.GatewayTrace)
		if strings.TrimSpace(traceId) == "" {
			traceId = constant.TraceIdPrefix + library.UniqueId()
		}
		c.Set(constant.TraceID, traceId)
	}else {
		traceId = traceIdInterface.(string)
	}

	return traceId
}

func GetRequestBody(c *gin.Context) interface{} {

	switch c.Request.Method {
	case http.MethodGet:
		return c.Request.URL.Query()

	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodPatch:
		var bodyBytes []byte // 我们需要的body内容

		// 从原有Request.Body读取
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println(err)
		}

		// 新建缓冲区并替换原有Request.body
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		var params interface{}
		json.Unmarshal(bodyBytes, &params)
		return params

	}

	return nil
}