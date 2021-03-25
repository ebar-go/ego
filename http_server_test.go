package ego

import (
	"fmt"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/egu"
	"github.com/gin-gonic/gin"
	"testing"
)

type IUserService interface {

}

type userService struct {
	userRepo struct{}
}

func NewUserService(userRepo struct{}) IUserService {
	return &userService{userRepo: userRepo}
}

type IUserHandler interface {
	Index(ctx *gin.Context)
}

type userHandler struct {
	userService IUserService
}

func NewUserHandler(userService IUserService) IUserHandler {
	return &userHandler{userService: userService}
}
func (userHandler) Index(ctx *gin.Context) {

}

func (handler userHandler) Hello(ctx *gin.Context) {
	fmt.Print(handler.userService)
	response.WrapContext(ctx).Success(nil)
}

func TestNewServer(t *testing.T) {
	app := App()

	egu.SecurePanic(app.LoadConfig("/usr/local/app.yaml"))

	err := app.Container().Invoke(func(router *gin.Engine, userHandler IUserHandler) {
		router.GET("index", userHandler.Index)
	})
	egu.SecurePanic(err)

	egu.SecurePanic(app.ServeHttp())

}
