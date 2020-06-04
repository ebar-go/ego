package http

import (
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/ego/utils/secure"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := NewServer()

	s.Router.Use(middleware.Favicon, middleware.RequestLog, middleware.Recover)
	s.Router.GET("/check", func(context *gin.Context) {
		response.WrapContext(context).Success("hello")
	})

	s.Router.GET("/error", func(context *gin.Context) {
		panic("err")
	})
	//_, fileStr, _, _ := runtime.Caller(0)
	//
	//utils.FatalError("ReadFromFile", config.ReadFromFile(filepath.Dir(fileStr) + "/../config/app.yaml"))

	//_ = event.DefaultDispatcher().Trigger(app.RedisConnectEvent, nil)
	secure.FatalError("StartHttpServer", s.Start())
}
