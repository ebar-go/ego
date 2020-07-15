package middleware

import (
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/egu"
	"github.com/gin-gonic/gin"
)

// Recover
func Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(*errors.Error)
			if ok {
				response.WrapContext(ctx).Error(err.Code, err.Message)

			} else {

				if app.Config().Server().Debug {
					app.Logger().Error("system_error", log.Context{
						"error": r,
						"trace": egu.RuntimeCaller(),
					})
				}

				response.WrapContext(ctx).Error(500, "System Error")
			}

		}
	}()
	ctx.Next()
}
