package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ebar-go/ego/http/response"
	"fmt"
	"github.com/ebar-go/ego/http/constant"
)

// NotFoundHandler 404
func NotFoundHandler(context *gin.Context)  {
	response.Error(context, constant.StatusNotFound, fmt.Sprintf("404 Not Found: %s", context.Request.RequestURI))
}
