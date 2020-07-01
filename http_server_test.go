package ego

import (
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/ego/utils/secure"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := HttpServer()

	s.Router.Use(middleware.Favicon, middleware.RequestLog, middleware.Recover)
	s.Router.Any("/check", func(context *gin.Context) {
		response.WrapContext(context).Success("hello")
	})


	secure.FatalError("StartHttpServer", s.Start(8080))
}
