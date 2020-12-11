package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zutim/ego/app"
	"github.com/zutim/ego/component/log"
	"github.com/zutim/ego/http/response"
	"github.com/zutim/egu"
)

// Recover
func Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			app.Logger().Debug("recover", log.Context{
				"error": r,
				"trace": egu.RuntimeCaller(),
			})

			err, ok := r.(*errors.Error)
			if ok {
				response.WrapContext(ctx).Error(err.Code, err.Message)

			} else {
				response.WrapContext(ctx).Error(500, "System Error")
			}

		}
	}()
	ctx.Next()
}
