package http

import (
	"fmt"
	"github.com/ebar-go/ego/http/handler"
	"github.com/gin-gonic/gin"
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

			server.Router.GET("/", func(context *gin.Context) {
				fmt.Println("hello,world")
			})
			if err := server.Start(tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
