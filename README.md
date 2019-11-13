[TOC]
## 说明
通用的GOLANG公共库,拒绝困扰程序员，提高生产力。

## 特性
- 提供简单的HTTP服务启动方式
- 提供简单的日志管理器
- 中间件(记录请求日志、JWT认证)
- 集成Apollo,支持热更新
- 集成Consul,支持服务注册、发现、注销
- 连接数据库,Redis

## 安装

```
go get -u github.com/ebar-go/ego
```

## 目录结构
```
.
├── cache
├── config
├── consul
├── db
├── http
├── library
├── log
├── task
├── test
├── go.mod
├── go.sum
└── README.md

```

## 模块
### http
基于gin框架的http服务器模块
- 启动web服务

```go
package main
import (
	"github.com/ebar-go/ego/http"
	"github.com/gin-gonic/gin"
	"fmt"
	)
func main() {
    server := http.Server {
        Address : "127.0.0.1", // 可以读取apollo地址
        Port:8088,
    }
    err := server.Init()
    if err != nil {
    	panic(err)
    }
    // TODO 添加路由
    server.Router.GET("/test", func(context *gin.Context) {
        fmt.Println("hello,world")
    })
    err =server.Start()
}
```

#### 中间件
- JWT
```go
package main
import (
	"github.com/ebar-go/ego/http"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/request"
	)
func main() {
    server := http.Server {
        Address : "127.0.0.1", // 可以读取apollo地址
        Port:8088,
    }
    err := server.Init()
    if err != nil {
    	panic(err)
    }
    // TODO 添加路由
    server.Router.GET("/test", func(context *gin.Context) {
        fmt.Println("hello,world")
    })
    
    api := server.Router.Group("/api")
    api.Use(middleware.JWT)
    {
    	api.GET("/user", func(context *gin.Context) {
    		
    	    fmt.Println("获取用户信息")
    	    fmt.Println(middleware.GetCurrentUser(context))
    	})
    }
    err =server.Start()
}
```

### cache
缓存模块,待定

### config
配置项,集成Apollo配置
- Apollo

```go
package main
import (
	"github.com/ebar-go/ego/config"
    "os"
	"fmt"
)
func main() {
    apollo := config.Apollo{
    	AppId: "open-api",
    	Cluster: "local",
    	Ip: "192.168.0.19:8080",
    	Namespace: "application",
    }
    if err := apollo.Init(); err != nil {
        // TODO 如果apollo启动失败，应该有备用方案
        fmt.Println("启动apollo失败:"+ err.Error())
        os.Exit(-1)
    }
    // 获取配置
    logFilePath := apollo.GetStringValue("LOG_FILE","/var/tmp")
    fmt.Println(logFilePath)
    // 另外，可以使用定时任务，监听配置变更
}

```
更多方法请查看测试用例

### consul
微服务(SOA),集成consul组件
- 服务注册

```go
package main
import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/ebar-go/ego/consul"
	"fmt"
	"github.com/ebar-go/ego/library"
)
func main() {
    config := consul.DefaultConfig()
    // 指定consul地址
    config.Address = "192.168.0.222:8500"
    client := &consul.Client{
        Config:config,
    }
    // 获取本机IP
    ip, err := library.GetLocalIp()  
    if err != nil {
    	panic(err)
    }
    registration := consul.NewServiceRegistration()
    registration.ID = "epet-go-demo-2"
    registration.Name = "project-demo"
    registration.Port = 8088
    registration.Tags = []string{"project-demo"}
    registration.Address = ip
    check := consul.NewServiceCheck()
    // 指定服务检查url
    check.HTTP = fmt.Sprintf("http://%s:%d%s", registration.Address, registration.Port, "/check")
    check.Timeout = "3s"
    check.Interval = "3s"
    check.DeregisterCriticalServiceAfter = "30s" //check失败后30秒删除本服务
    registration.Check = check
    err = client.Register(registration)
}
```
- 服务发现

```go
package main
import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/ebar-go/ego/consul"
	"fmt"
	"github.com/ebar-go/ego/library"
)
func main() {
    config := consul.DefaultConfig()
    // 指定consul地址
    config.Address = "192.168.0.222:8500"
    client := &consul.Client{
        Config:config,
    }
    items, err := client.Discover("project-demo")
    if err != nil {
    	panic(err)
    }
    service, err := client.LoadBalance(items)
    fmt.Println(service.GetHost())
}
```

更多方法请查看测试用例

### library
公共库
- Date

```go
package main
import ("github.com/ebar-go/ego/library"
    "fmt"
)
func main() {//
    // 获取当前时间
    now := library.GetTimeStr()
    fmt.Println(now)
}

```

- 打印调试

```go
package main
import (
	"github.com/ebar-go/ego/library"
)
func main() {
	// 打印调试
    library.Debug(123, "aaa")
}

```


### log
日志管理器
- 输出到控制台

```go
package main
import (
       	"github.com/ebar-go/ego/log"
       	"os"
       )
func main() {
    logger := log.New()
    logger.Debug("test debug", log.Context{"name":"123"})
}
```

- 输出到文件

```go
package main
import (
	"github.com/ebar-go/ego/log"
	"os"
	"fmt"
)
func main() {
    logger := log.New()
    filePath := "/var/log/system.log"
    fmt.Println(filePath)
    file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err == nil {
    	logger.SetOutWriter(file)
    }else {
    	fmt.Println("err:" + err.Error())
	}
    logger.Debug("test debug", log.Context{"name":"123"})
}
```

### HTTP请求
- kong
- http

更多方法请查看测试用例

### mysql
数据库

### redis
redis

### test 单元测试

## TODO
- 参数验证器