package request

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Request struct {
	context *gin.Context
}

func Default(context *gin.Context) *Request {
	return &Request{
		context: context,
	}
}

// GetDefaultQueryInt 获取int参数,支持指定默认值
func (req *Request) GetDefaultQueryInt(key string, defaultValue int) int{
	param := req.context.Query(key)
	result , err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}

	return result
}

// GetQueryInt 获取int参数
func (req *Request) GetQueryInt(key string, defaultValue int) (int, error){
	param := req.context.Query(key)
	result , err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	return result, nil
}
