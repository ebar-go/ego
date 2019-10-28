package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ebar-go/ego/http/response"
	"fmt"
)

// NotFoundHandler 404
func NotFoundHandler(context *gin.Context)  {
	responseWriter := response.Default(context)

	responseWriter.Error(404, fmt.Sprintf("404 Not Found: %s", context.Request.RequestURI))
}
