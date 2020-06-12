package app

import (
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/log"
	"net"
	"net/http"
	"time"
)


// InitDB 初始化DB
func InitDB() error {
	//dialect := "mysql"
	//group := config.MysqlGroup()
	//if group.Items == nil {
	//	return fmt.Errorf("mysql config is empty")
	//}
	//
	//for name, item := range group.Items {
	//	dataSourceItems := item.DsnItems()
	//
	//	adapter, err := mysql.NewReadWriteAdapter(dialect, dataSourceItems)
	//	if err != nil {
	//		return err
	//	}
	//
	//	adapter.SetMaxIdleConns(item.MaxIdleConnections)
	//	adapter.SetMaxOpenConns(item.MaxOpenConnections)
	//	adapter.SetConnMaxLifetime(time.Duration(item.MaxLifeTime) * time.Second)
	//
	//	conn, err := gorm.Open(dialect, adapter)
	//	if err != nil {
	//		return err
	//	}
	//	dbGroup[name] = conn
	//}
	//
	//event.Trigger(event.AfterDatabaseConnect, nil)

	return nil
}

// newHttpClient
func newHttpClient(conf *config.Config) *http.Client {
	return &http.Client{
		Transport: &http.Transport{ // 配置连接池
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			IdleConnTimeout: time.Duration(conf.Server().HttpRequestTimeOut) * time.Second,
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Duration(conf.Server().HttpRequestTimeOut) * time.Second,
	}
}

// newLogger
func newLogger (conf *config.Config) *log.Logger {
	return log.New(conf.Server().LogPath,
		conf.Server().Debug,
		map[string]interface{}{
			"system_name": conf.Server().Name,
		})
}