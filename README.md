## 说明
通用的GOLANG公共库,拒绝困扰程序员，提高生产力。

## 特性
- 提供简单的HTTP服务启动方式
- 提供简单的日志管理器
- 自由的集成组件,拒绝臃肿
- 中间件(记录请求日志、JWT认证)
- 集成Apollo,支持热更新
- 集成Consul,支持服务注册、发现、注销
- Mysql,Redis连接池
- 集成prometheus监控
- 集成MNS
- 集成参数验证器
- 集成websocket
- 集成配置文件,支持json,ini,yaml等格式
- 集成di,统一管理全局变量

## 安装

```
go get -u github.com/ebar-go/ego
```

## 模块
### http
基于gin框架的http服务器模块
- 启动web服务

```go
package main
import (

"fmt"
"github.com/ebar-go/ego/helper"
"github.com/ebar-go/ego/http"
"github.com/gin-gonic/gin"
	)

func main() {
    server := http.NewServer()
    // 添加路由
    server.Router.GET("/test", func(context *gin.Context) {
        fmt.Println("hello,world")
    })
    
    helper.FatalError("StartServer", server.Start())
}
```

#### 中间件
集成JWT,请求日志,跨域等中间件
- JWT
```go
package main
import (
	"github.com/ebar-go/ego/http"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/helper"
	)
func main() {
    server := http.NewServer()
    // 添加路由
    server.Router.GET("/test", func(context *gin.Context) {
        fmt.Println("hello,world")
    })
    
    api := server.Router.Group("/api")
    api.Use(middleware.JWT)
    {
    	api.GET("/user", func(context *gin.Context) {
    		
    	    fmt.Println("获取用户信息")
    	    fmt.Println(middleware.GetCurrentClaims(context))
    	})
    }
    helper.FatalError("StartServer", server.Start())
}
```

- RequestLog 请求日志

- Cors
```go
package main
import (
"github.com/ebar-go/ego/http"
"github.com/ebar-go/ego/http/middleware"
)
func main() {
    server := http.NewServer()
    // 根据情况，中间件一般放在路由的前面,否则会在handler处理完成后生效
    server.Router.Use(middleware.CORS) // 跨域
    server.Router.Use(middleware.Recover) // 错误处理机制
}
```

### DI 全局变量管理
通过集成`https://github.com/uber-go/dig`,用依赖注入的方式管理全局变量
```go
package main
import (

"fmt"
"github.com/ebar-go/ego/app"
"github.com/ebar-go/ego/helper"
)
func main() {
    // 系统容器,将配置、数据库、Redis连接、mns客户端等
    fmt.Println(app.Config())
    
    // 自定义容器
    container := app.NewContainer()
    container.Provide(func() int{
        globalV := 1
        return globalV
    })
    err := container.Invoke(func(globalV int) {
        fmt.Println(globalV)
    })
    helper.CheckError("GetGlobalV", err)

}
```

### cache
推荐使用 go-cache

### config
配置项,集成Apollo配置
- Apollo

```go
package main
import (
	"github.com/ebar-go/ego/component/apollo"
    "os"
	"fmt"
	"github.com/ebar-go/ego/helper"
)
func main() {
    conf := apollo.Conf{
    	AppId: "open-api",
    	Cluster: "local",
    	Ip: "192.168.0.19:8080",
    	Namespace: "application",
    }
    helper.FatalError("InitApollo", apollo.Init(conf))
    
    // 获取配置
    logFilePath := apollo.GetStringValue("LOG_FILE","/var/tmp")
    fmt.Println(logFilePath)

    // 另外，可以使用定时任务，监听配置变更,更多请查看test
}

```

### consul
微服务(SOA),集成consul组件
- 服务注册

```go
package main
import (
	"github.com/ebar-go/ego/component/consul"
	"fmt"
	"github.com/ebar-go/ego/helper"
)
func main() {
    config := consul.DefaultConfig()
    // 指定consul地址
    config.Address = "192.168.0.222:8500"
    client := &consul.Client{
        Config:config,
    }
    // 获取本机IP
    ip, err := helper.GetLocalIp()  
    if err != nil {
    	panic(err)
    }
    registration := consul.NewServiceRegistration()
    registration.ID = "go-demo-2"
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
	"github.com/ebar-go/ego/component/consul"
	"fmt"
	"github.com/ebar-go/ego/helper"
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
    if err != nil {
        panic(err)
    }   
    fmt.Println(service.GetHost())
}
```

更多方法请查看测试用例

### helper
公共库

```go
package main
import (
    "github.com/ebar-go/ego/helper"
    "fmt"
)
func main() {//
    // 获取当前时间
    fmt.Println("获取当前时间:" , helper.GetTimeStr())
    helper.Debug("打印调试")
    
    // 更多请查看使用示例
}
```

### log
日志管理器
- 系统日志

```go
package main
import (
       	"github.com/ebar-go/ego/app"
        "github.com/ebar-go/ego/component/log"
       	"os"
       )
func main() {
    app.LogManager().App().Info("test", log.Context{"a":1})
    app.LogManager().App().Debug("test", log.Context{"a":1})
    app.LogManager().App().Warn("test", log.Context{"a":1})
    app.LogManager().App().Error("test", log.Context{"a":1})
}
```
- 自定义

```go
package main
import (
       	"github.com/ebar-go/ego/component/log"
       	"os"
       )
func main() {
    logger := log.New()
    logger.Debug("test debug", log.Context{"name":"123"})
}
```

### HTTP请求客户端
提供官方的http包、fasthttp(推荐使用),kong网关的http客户端
```go
package main
import (
       	"github.com/ebar-go/ego/component/mysql"
       	"github.com/ebar-go/ego/http/client"
       	"github.com/ebar-go/ego/http/client/request"
       	"os"
       	"fmt"
       )
func main() {
	// 官方http,支持长连接
	httpClient := client.NewHttpClient()
	
	// fasthttp
	// httpClient := client.NewFastHttpClient()
	
	// kong
	// httpClient := client.NewKongClient()
	// 设置其他参数...
	req := httpClient.NewRequest(request.Param{
		Url: "http://localhost:8080/test",
		Method: request.Get,
	})
	
	resp, err := httpClient.Execute(req)
	fmt.Println(resp, err)
}
```

### mysql
集成的Gorm,使用连接池

```go
package main
import (
       	"github.com/ebar-go/ego/app"
       )
func main() {
    conn := app.Mysql()
    defer conn.Close()
    
    conn.DB().Ping()
    
}
```

### redis
集成的go-redis,使用连接池

```go
package main
import (
       	"github.com/ebar-go/ego/app"
       	"github.com/ebar-go/ego/helper"
       	"os"
       	"fmt"
       )
func main() {
    conn := app.Redis()
    if err := conn.Set("key", "value", 0).Err(); err != nil {
    	fmt.Println(err)
    }
    
    val, err := conn.Get("key").Result()
    fmt.Println("key", val, err)
}
```

### 对接阿里云MNS
```go
package main
import (
       	"github.com/ebar-go/ego/app"
"github.com/ebar-go/ego/component/mns"
       	"github.com/ebar-go/ego/helper"
       	"os"
       	"fmt"
       )
func main()  {
	// 添加队列处理方法
	app.Mns().AddQueue("queueName",  Process, 30)
    
    // 使用定时任务监听队列
    _ = app.Task().AddFunc(fmt.Sprintf("*/%d * * * * *", 1), func() {
    	app.Mns().ListenQueues()
    })
    
    // 通过主体发送消息
    resp, err := app.Mns().PublishMessage("topicName", mns.Params{}, "filterTag")
    fmt.Println(resp, err)
}

func Process(params mns.Params) error {
    fmt.Println(params)
    return nil
}
```

### 对接prometheus监控
监控Mysql

```go
package main
import (
       	"github.com/ebar-go/ego/app"
"github.com/ebar-go/ego/component/mysql"
       	"github.com/ebar-go/ego/component/prometheus"
       	"github.com/ebar-go/ego/helper"
       	"github.com/ebar-go/ego/http"
       	"os"
       )
func main() {
    
    conn := app.Mysql()
    prometheus.ListenMysql(conn, "server")
    
    server := http.NewServer()
    
    server.Router.GET("/metrics", prometheus.Handler)
    
    helper.FatalError("StartServer", server.Start())
    
}
```

### 参数验证器
更多验证规则请查阅: [https://github.com/go-playground/validator](https://github.com/go-playground/validator)，经过社区验证的高可用、高性能验证器
```go
package main
import (
       	"github.com/ebar-go/ego/helper"
       	"github.com/ebar-go/ego/http"
       	"github.com/ebar-go/ego/http/response"
       	"os"
       	"github.com/gin-gonic/gin"
       	"github.com/gin-gonic/gin/binding"
       	"github.com/ebar-go/ego/http/validator"
       )

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func main() {
	// 使用自定义验证器代替gin的validator
	binding.Validator = new(validator.Validator)
	
    server := http.NewServer()
    
    server.Router.GET("/login", func(c *gin.Context) {
        var json Login
        if err := c.ShouldBindJSON(&json); err != nil {
        	// 参数错误的业务码自定义
        	invalidParamCode := 1001
        	response.Error(c, invalidParamCode, err.Error())
        	return
        }
        		
       // .. 其他代码
        		
       response.Success(c, response.Data{
       	"token":"1234567abcd",
       })
    })
    
    helper.FatalError("StartServer", server.Start())
    
}
```

### test 单元测试
```go
package test
import (

"github.com/magiconair/properties/assert"
"testing"
)

func TestName(t *testing.T) {
    expect := 1
    got := 1
    assert.Equal(t, expect, got)
}
```

### 支持websocket
基于`github.com/gorilla/websocket`实现websocket

```go
package main
import (
	"github.com/ebar-go/ego/app"
"github.com/ebar-go/ego/http"
	"github.com/ebar-go/ego/ws"
	"github.com/gin-gonic/gin"
	"fmt"
    "github.com/ebar-go/ego/helper"
	nethttp "net/http"
	)

func main() {
	// 使用协程启动
	go app.WebSocket().Start()
	
    server := http.NewServer()
    // 添加路由
    server.Router.GET("/test", func(context *gin.Context) {
        fmt.Println("hello,world")
    })
    server.Router.GET("/ws", func(context *gin.Context) {
        conn, err := ws.GetUpgradeConnection(context.Writer, context.Request)
        
        if err != nil {
        	nethttp.NotFound(context.Writer, context.Request)
        	return
        }
        
        
        // TODO 根据tag扩展handler
        
        client := ws.NewClient(conn, func(ctx *ws.Context) string {
        	// do something
            return ctx.GetMessage()
        })
        
        // 将用户信息放入扩展里
        client.Extends["user"] = struct {
         Name string
        }{Name: "test"}
        
        app.WebSocket().RegisterClient(client)
        
        go client.Listen()
    })
    
    helper.FatalError("StartServer", server.Start())
}
```
## TODO
- 支持RPC
