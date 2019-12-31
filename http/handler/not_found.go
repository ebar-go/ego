package handler

import (
	"fmt"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/http/response"
	"github.com/gin-gonic/gin"
)

// NotFoundHandler 404
func NotFoundHandler(ctx *gin.Context) {
	response.WrapperContext(ctx).Error(errors.NotFound())
	response.WrapperContext(ctx).Error(404, fmt.Sprintf("404 Not Found: %s", ctx.Request.RequestURI))
}
