package middleware

import (
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/ego/utils"
	"github.com/gin-gonic/gin"
)

// Recover
func Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(*errors.Error)
			if ok {
				response.WrapContext(ctx).Error(err.Code, err.Message)

			}else {
				log.Error("system_error", log.Context{
					"error": r,
					"trace": utils.Trace(),
				})
				response.WrapContext(ctx).Error(500, "System Error")
			}



		}
	}()
	ctx.Next()
}
