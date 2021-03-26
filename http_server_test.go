package ego

import (
	"fmt"
	"github.com/ebar-go/ego/component/auth"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/mysql"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/egu"
	"github.com/gin-gonic/gin"
	"testing"
)

type IUserService interface {
}

type IUserRepo interface {
}

type userRepo struct {
	db mysql.Database
}

func NewUserRepo(db mysql.Database) IUserRepo {
	return &userRepo{db}
}

type userService struct {
	userRepo IUserRepo
}

func NewUserService(userRepo IUserRepo) IUserService {
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

	app.Container().Provide(NewUserRepo)
	app.Container().Provide(NewUserService)
	app.Container().Provide(NewUserHandler)

	err := app.Container().Invoke(func(router *gin.Engine,
		logger *log.Logger,
		jwtAuth auth.Jwt,
		userHandler IUserHandler,
	) {
		router.Use(middleware.CORS, middleware.Recover(logger))
		router.GET("index", userHandler.Index)
		router.GET("home", middleware.JWT(jwtAuth))

	})
	egu.SecurePanic(err)

	egu.SecurePanic(app.ServeHttp())

}
