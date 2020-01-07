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
#### 服务初始化
基于gin框架的http服务器模块
- 启动web服务

```go
package main
import (
    "fmt"
    "github.com/ebar-go/ego/http"
    "github.com/ebar-go/ego/utils"
    "github.com/gin-gonic/gin"
)

func main() {
    server := http.NewServer()
    // 添加路由
    server.Router.GET("/test", func(context *gin.Context) {
        fmt.Println("hello,world")
    })
    
    utils.FatalError("StartServer", server.Start())
}
```

#### http请求客户端
支持原生http、fasthttp、kong网关等
```go
package main
import (
    "fmt"
    "github.com/ebar-go/ego/http/client"
    "github.com/ebar-go/ego/http/client/request"
)

func main() {
    cli := client.NewFastHttpClient()
    // cli := client.NewHttpClient()
    // cli := client.NewKongClient()
    req := cli.NewRequest(request.Param{
        Method:"GET",
        Url:"http://localhost:8080/index",
        Headers:nil,
        Body:nil,
    })
    resp, err := cli.Execute(req)
    fmt.Println(resp, err)

}
```

#### 中间件
- CORS
跨域中间件

- JWT
JsonWebToken鉴权

- RequestLog
请求日志记录

- Recover
处理panic的中间件
```go
package main
import (
    "fmt"
    "github.com/ebar-go/ego/http"
    "github.com/ebar-go/ego/http/middleware"
    "github.com/ebar-go/ego/utils"
    "github.com/gin-gonic/gin"
)

func main() {
    server := http.NewServer()
    // middleware must used before init router 
    // Recover middleware
    server.Router.Use(middleware.Recover)
    // JWT middleware,make sure Config.JwtKey is not empty
    server.Router.Use(middleware.JWT)
    // CORS middleware
    server.Router.Use(middleware.CORS)
    // RequestLog middleware,not necessary, system auto load
    server.Router.Use(middleware.RequestLog)

    // Add router
    server.Router.GET("/test", func(context *gin.Context) {
        fmt.Println("hello,world")
    })
    
    utils.FatalError("StartServer", server.Start())
}
```


#### 分页组件
对数据进行分页

#### 响应
固定格式的json输出
```go
package handler
import (
    "github.com/ebar-go/ego/http/response"
    "github.com/gin-gonic/gin"
)
func IndexHandler(ctx *gin.Context) {
    response.WrapContext(ctx).Success(response.Data{"hello":"world"})
    // response.WrapContext(ctx).Error(500, "system error")
}
```

#### validator
集成`github.com/go-playground/validator`

### 组件包
#### apollo
```go
err := apollo.Init(apollo.Conf{
    AppId:            "",
    Cluster:          "",
    Namespace:        "",
    Ip:               "",
    BackupConfigPath: "",
})
```
#### auth
#### consul
微服务(SOA),集成consul组件
- 服务注册
- 服务发现

#### 日志
#### mns
阿里云mns客户端
```go
// trigger mns init event
// make sure config is not empty 
app.EventDispather().Trigger(app.MNSClientInitEvent, nil)
mns := app.Mns()
```

#### mysql
- connect database
```go
// trigger mysql connect event
// make sure config is not empty 
app.EventDispather().Trigger(app.MysqlConnectEvent, nil)
db := app.Mysql() 
// .. curd
```
- 自定义model

#### prometheus
集成prometheus

#### trace
全局ID
```go
package main
import (
 "fmt"
 "github.com/ebar-go/ego/component/trace"
)
func main() {
    // set global trace id first
    trace.SetTraceId(trace.NewId())
    // use defer execute recycle
    defer trace.DeleteTraceId()
    // get id in other location
    fmt.Println(trace.GetTraceId())
}
```

### 配置项
读取配置项

### utils
常用工具库
```go
utils.FatalError("Failt error", func () error{
    return errors.New("test")
})
utils.LogError("Log error", nil)
```
- 数组
- 转换
- 日期
- 文件
- json

提供json操作
```go
json.Encode(someValue)
json.Decode(jsonStr, &someValue)
```
- 数字

提供数字型相关常用方法
```go
// default int 
httpPort = number.DefaultInt(httpPort, 8080)
```
- 字符串

提供字符串相关常用方法
```go
// default string
host = strings.Default(host, "127.0.0.1")
```

### log
日志管理器
- 系统日志
```go
package main
import (
 "fmt"
 "github.com/ebar-go/ego/app"
"github.com/ebar-go/ego/component/log"
)
func main() {
    // if not run http server,please trigger log manager init event
    app.EventDispatcher().Trigger(app.LogManagerInitEvent, nil)
    app.LogManager().App().Info("test", log.Context{"hello":"world"})
    app.LogManager().App().Debug("debug", log.Context{"hello":"world"})
    app.LogManager().App().Warn("warn", log.Context{"hello":"world"})
    app.LogManager().App().Error("error", log.Context{"hello":"world"})
}

```
