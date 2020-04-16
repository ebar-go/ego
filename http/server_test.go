package http

import (
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/ego/utils"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := NewServer()

	s.Router.Use( middleware.Favicon, middleware.RequestLog, middleware.Recover)
	s.Router.GET("/", func(context *gin.Context) {
		response.WrapContext(context).Success("hello")
	})

	s.Router.GET("/error", func(context *gin.Context) {
		panic("err")
	})
	_, fileStr, _, _ := runtime.Caller(0)

	utils.FatalError("ReadFromFile", config.ReadFromFile(filepath.Dir(fileStr) + "/../config/app.yaml"))

	//_ = event.DefaultDispatcher().Trigger(app.RedisConnectEvent, nil)
	utils.FatalError("StartHttpServer",s.Start())
}
