package app

import (
	"github.com/ebar-go/ego/component/auth"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/mns"
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ego/event"
	"github.com/ebar-go/ego/ws"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron"
	"go.uber.org/dig"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	EventDispatcher().Trigger(ConfigInitEvent, nil)
	m.Run()
}

func TestConfig(t *testing.T) {
	tests := []struct {
		name     string
		wantConf *config.Config
	}{
		{
			name:     "get",
			wantConf: config.NewInstance(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotConf := Config(); !reflect.DeepEqual(gotConf, tt.wantConf) {
				t.Errorf("Config() = %v, want %v", gotConf, tt.wantConf)
			}
		})
	}
}

func TestEventDispatcher(t *testing.T) {
	tests := []struct {
		name           string
		wantDispatcher event.Dispatcher
	}{
		{
			name:           "test",
			wantDispatcher: event.NewDispatcher(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDispatcher := EventDispatcher(); !reflect.DeepEqual(gotDispatcher, tt.wantDispatcher) {
				t.Errorf("EventDispatcher() = %v, want %v", gotDispatcher, tt.wantDispatcher)
			}
		})
	}
}

func TestJwt(t *testing.T) {
	tests := []struct {
		name    string
		wantJwt auth.Jwt
	}{
		{
			name:    "key",
			wantJwt: &auth.JwtAuth{SignKey:[]byte("key")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Config().JwtSignKey = []byte(tt.name)
			EventDispatcher().Trigger(JwtInitEvent, nil)
			if gotJwt := Jwt(); !reflect.DeepEqual(gotJwt, tt.wantJwt) {
				t.Errorf("Jwt() = %v, want %v", gotJwt, tt.wantJwt)
			}
		})
	}
}

func TestLogManager(t *testing.T) {
	tests := []struct {
		name        string
		wantManager log.Manager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotManager := LogManager(); !reflect.DeepEqual(gotManager, tt.wantManager) {
				t.Errorf("LogManager() = %v, want %v", gotManager, tt.wantManager)
			}
		})
	}
}

func TestMns(t *testing.T) {
	tests := []struct {
		name       string
		wantClient mns.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotClient := Mns(); !reflect.DeepEqual(gotClient, tt.wantClient) {
				t.Errorf("Mns() = %v, want %v", gotClient, tt.wantClient)
			}
		})
	}
}

func TestMysql(t *testing.T) {
	tests := []struct {
		name           string
		wantConnection *gorm.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotConnection := Mysql(); !reflect.DeepEqual(gotConnection, tt.wantConnection) {
				t.Errorf("Mysql() = %v, want %v", gotConnection, tt.wantConnection)
			}
		})
	}
}

func TestNewContainer(t *testing.T) {
	tests := []struct {
		name string
		want *dig.Container
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContainer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis(t *testing.T) {
	tests := []struct {
		name           string
		wantConnection *redis.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotConnection := Redis(); !reflect.DeepEqual(gotConnection, tt.wantConnection) {
				t.Errorf("Redis() = %v, want %v", gotConnection, tt.wantConnection)
			}
		})
	}
}

func TestTask(t *testing.T) {
	tests := []struct {
		name        string
		wantManager *cron.Cron
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotManager := Task(); !reflect.DeepEqual(gotManager, tt.wantManager) {
				t.Errorf("Task() = %v, want %v", gotManager, tt.wantManager)
			}
		})
	}
}

func TestWebSocket(t *testing.T) {
	tests := []struct {
		name        string
		wantManager ws.Manager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotManager := WebSocket(); !reflect.DeepEqual(gotManager, tt.wantManager) {
				t.Errorf("WebSocket() = %v, want %v", gotManager, tt.wantManager)
			}
		})
	}
}

func Test_connectDatabase(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := connectDatabase(); (err != nil) != tt.wantErr {
				t.Errorf("connectDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_connectRedis(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := connectRedis(); (err != nil) != tt.wantErr {
				t.Errorf("connectRedis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initConfig(); (err != nil) != tt.wantErr {
				t.Errorf("initConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initEventDispatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initEventDispatcher(); (err != nil) != tt.wantErr {
				t.Errorf("initEventDispatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initJwt(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initJwt(); (err != nil) != tt.wantErr {
				t.Errorf("initJwt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initLogManager(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initLogManager(); (err != nil) != tt.wantErr {
				t.Errorf("initLogManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initMnsClient(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initMnsClient(); (err != nil) != tt.wantErr {
				t.Errorf("initMnsClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initTaskManager(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initTaskManager(); (err != nil) != tt.wantErr {
				t.Errorf("initTaskManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initWebSocketManager(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initWebSocketManager(); (err != nil) != tt.wantErr {
				t.Errorf("initWebSocketManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}