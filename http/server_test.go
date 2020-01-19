package http

import (
	"fmt"
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/ego/component/consul"
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ego/http/handler"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/ego/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
		want *Server
	}{
		{
			name: "test",
			want: &Server{
				mu:              sync.Mutex{},
				Router:          gin.Default(),
				NotFoundHandler: handler.NotFoundHandler,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(); got == nil {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	type fields struct {
		mu              sync.Mutex
		Router          *gin.Engine
		NotFoundHandler func(ctx *gin.Context)
	}
	type args struct {
		args []int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test",
			fields:  fields{
				mu:              sync.Mutex{},
				Router:          gin.Default(),
				NotFoundHandler: nil,
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{
				mu:              tt.fields.mu,
				Router:          tt.fields.Router,
				NotFoundHandler: tt.fields.NotFoundHandler,
			}
			server.Router.Use(middleware.RequestLog)

			server.Router.GET("/", func(context *gin.Context) {
				fmt.Println("hello,world")
			})

			server.Router.POST("/post", func(context *gin.Context) {
				fmt.Println("hello,world")
			})
			if err := server.Start(tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getClient() *consul.Client {
	config := consul.DefaultConfig()
	config.Address = "192.168.0.222:8500"

	return &consul.Client{
		Config: config,
	}
}

func TestServerRegister(t *testing.T)  {
	server := NewServer()
	server.Router.Use(middleware.RequestLog)

	server.Router.GET("/check", func(context *gin.Context) {
		fmt.Println("hello,world")
		response.WrapContext(context).Success(nil)
	})

	server.Setup()
	config.Server().Port = 8080

	client := getClient()

	ip, err := utils.GetLocalIp()
	assert.Nil(t, err)

	registration := consul.NewServiceRegistration()
	registration.ID = "epet-go-demo-1"
	registration.Name = "epet-go-demo"
	registration.Port = config.Server().Port
	registration.Tags = []string{"epet-go-demo"}
	registration.Address = ip


	check := consul.NewServiceCheck()
	check.HTTP = fmt.Sprintf("http://%s:%d%s", registration.Address, registration.Port, "/check")
	check.Timeout = "3s"
	check.Interval = "3s"
	check.DeregisterCriticalServiceAfter = "30s" //check失败后30秒删除本服务
	registration.Check = check

	err = client.Register(registration)

	utils.SecurePanic(server.Start())
}

func TestServerRegister2(t *testing.T)  {
	server := NewServer()
	server.Router.Use(middleware.RequestLog)

	server.Router.GET("/check", func(context *gin.Context) {
		fmt.Println("hello,world2")
		response.WrapContext(context).Success(nil)
	})

	server.Setup()
	config.Server().Port = 8081

	client := getClient()

	ip, err := utils.GetLocalIp()
	assert.Nil(t, err)

	registration := consul.NewServiceRegistration()
	registration.ID = "epet-go-demo-2"
	registration.Name = "epet-go-demo"
	registration.Port = app.Config().Server().Port
	registration.Tags = []string{"epet-go-demo"}
	registration.Address = ip
	registration.Weights.Warning = 2

	check := consul.NewServiceCheck()
	check.HTTP = fmt.Sprintf("http://%s:%d%s", registration.Address, registration.Port, "/check")
	check.Timeout = "3s"
	check.Interval = "3s"
	check.DeregisterCriticalServiceAfter = "30s" //check失败后30秒删除本服务
	registration.Check = check

	err = client.Register(registration)

	utils.SecurePanic(server.Start())
}
