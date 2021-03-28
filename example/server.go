package main

import (
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/component/auth"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/mysql"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/egu"
	"github.com/gin-gonic/gin"
)

func main() {
	app := ego.App()

	egu.SecurePanic(app.LoadConfig("./app.yaml"))
	//
	_ = app.Container().Provide(NewUserRepo)
	_ = app.Container().Provide(NewUserService)
	_ = app.Container().Provide(NewUserHandler)

	err := app.Container().Invoke(func(router *gin.Engine,
		logger *log.Logger,
		jwtAuth auth.Jwt,
		userHandler IUserHandler,
	) {
		router.Use(middleware.CORS, middleware.Recover, middleware.RequestLog(logger))
		router.GET("index", userHandler.Index)
		router.GET("home", middleware.JWT(jwtAuth))
		router.GET("log", func(ctx *gin.Context) {
			logger.Info("test", log.Context{"test":"content"})
			response.WrapContext(ctx).Success(nil)
		})

	})
	egu.SecurePanic(err)

	app.ListenHTTP()
	app.ListenCron()

	app.Serve()
}


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
func (handler userHandler) Index(ctx *gin.Context) {
	response.WrapContext(ctx).Success(nil)
}

func (handler userHandler) Hello(ctx *gin.Context) {
	response.WrapContext(ctx).Success(response.Data{"hello":"world"})
}