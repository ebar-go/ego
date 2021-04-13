package middleware

import (
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/http/response"
	"github.com/gin-gonic/gin"
	"log"
	"runtime"
)

// Recover
func Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {

			err, ok := r.(*errors.Error)
			if ok {
				response.WrapContext(ctx).Error(err.Code, err.Message)
			} else {
				log.Println(r)
				log.Println(string(stack()))
				response.WrapContext(ctx).Error(500, "System Error")
			}

		}
	}()
	ctx.Next()

}

func stack() []byte {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	return buf
}
