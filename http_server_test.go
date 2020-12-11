package ego

import (
	"github.com/gin-gonic/gin"
	"github.com/zutim/ego/http/response"
	"github.com/zutim/egu"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := HttpServer()

	//s.Router.Use(middleware.Favicon, middleware.RequestLog, middleware.Recover)
	s.Router.GET("/list", func(context *gin.Context) {
		response.WrapContext(context).Success("hello")
	})

	egu.FatalError("StartHttpServer", s.Start(8081))
}
