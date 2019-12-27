package handler

import (
	"fmt"
	"github.com/ebar-go/ego/http/response"
	"github.com/gin-gonic/gin"
)

// NotFoundHandler 404
func NotFoundHandler(context *gin.Context) {
	response.Error(context, 404, fmt.Sprintf("404 Not Found: %s", context.Request.RequestURI))
}
