package request

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/ebar-go/ego/http/middleware"
)

type Request struct {
	context *gin.Context
}

func Default(context *gin.Context) *Request {
	return &Request{
		context: context,
	}
}

// GetCurrentUser 获取当前用户
func (req *Request) GetCurrentUser() *middleware.User {
	user, exist := req.context.Get(middleware.JwtUserKey)
	if !exist {
		return nil
	}

	return user.(*middleware.User)
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
func (req *Request) GetQueryInt(key string) (int, error){
	param := req.context.Query(key)
	result , err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	return result, nil
}
