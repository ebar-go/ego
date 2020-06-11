## 说明
基于Gin的web微服务框架。

## 特性
- 为分布式而生的http,websocket服务。
- 提供强大的日志管理器，支持按日期自动分割。
- 丰富的中间件(请求日志、JWT认证,跨域,Recover,全局链路)
- 支持基于consul,etcd的服务注册、发现、注销
- 基于gorm扩展的支持多库连接，读写分离的mysql组件
- redis组件
- 自定义参数验证器
- 配置文件,支持json,ini,yaml等格式

## 安装

```
go get -u github.com/ebar-go/ego
```

## 示例
完整的web项目请查看[github.com/ebar-go/ego-demo](https://github.com/ebar-go/ego-demo)

## 模块
### http服务
```go
package main
import (
"fmt"
"github.com/ebar-go/ego"
"github.com/gin-gonic/gin"
"github.com/ebar-go/ego/utils/secure"
"time"
)
func main() {
    server := ego.HttpServer()
    // 添加路由
    server.Router.GET("/test", func(context *gin.Context) {
        fmt.Println("hello,world")
    })
    secure.Panic(server.Start())
}
// 支持平滑重启
func init() {
    // 支持停止http服务时的回调
    event.Listen(event.BeforeHttpShutdown, func(ev event.Event) {
        // 让服务器等待3秒，等待异步业务执行结束
        time.Sleep(time.Second * 3)
    	// 关闭数据库
    	fmt.Println("close database")
    	_ = app.DB().Close()
    })
} 
```

### websocket服务
支持与http服务共存的websocket服务   
```go
package main
import (
"github.com/ebar-go/ego"
"github.com/ebar-go/ego/utils/secure"
"github.com/gin-gonic/gin"
"github.com/ebar-go/ego/http/response"
)
func main() {
    s := ego.HttpServer()
    ws := ego.WebsocketServer()

	s.Router.GET("/check", func(context *gin.Context) {
		response.WrapContext(context).Success("hello")
	})
	s.Router.GET("/ws", func(ctx *gin.Context) {
		// get websocket conn
		conn, err := ws.UpgradeConn(ctx.Writer, ctx.Request)
		if err != nil {
			secure.Panic(err.Error())
		}

		ws.Register(conn, func(message []byte){
			if string(message) == "broadcast" {// 广播
				ws.Broadcast([]byte("hello,welcome"), nil)
				return
			}
			ws.Send(message, conn) // 单对单发送

		})
	})

	go ws.Start()

	secure.FatalError("StartHttpServer", s.Start())
}
```

### 配置
集成[https://github.com/spf13/viper](https://github.com/spf13/viper)
```go
// 或者从文件读取
config.ReadFromFile(configFilePath)
// 读取配置
fmt.Println(viper.GetString("someKey"))
```

#### 发起http请求
```go
url := "http://baidu.com"
request,_ := http.NewRequest(http.MethodGet, url, nil)
resp, err := curl.Execute(request)
fmt.Println(resp, err)
```

#### 中间件

```go
server := http.NewServer()
// middleware must used before init router 
// Recover middleware
server.Router.Use(middleware.Recover)
// JWT middleware,make sure Config.JwtKey is not empty
server.Router.Use(middleware.JWT)
// CORS middleware
server.Router.Use(middleware.CORS)
// RequestLog middleware
server.Router.Use(middleware.RequestLog)
```

#### 响应
固定格式的json输出
```go
// 输出成功的响应数据
response.WrapContext(ctx).Success(response.Data{"hello":"world"})
// 数据错误的响应数据
response.WrapContext(ctx).Error(1001, "some error")
```

#### 数据校验器
基于`github.com/go-playground/validator`,支持自定义字段名称`comment`

- handler里校验参数
```go
type AuthRequest struct {
    // 如果是表单提交，使用form,否则获取不到数据
    Email string `json:"email" validate:"required,email" comment:"邮箱"` // 验证邮箱格式
    Pass string `json:"pass" binding:"required,min=6,max=10" comment:"密码"` // 验证密码，长度为6~10
}
    
var request AuthRequest
if err := ctx.ShouldBindJSON(&request); err != nil {
    secure.Panic(errors.New(1001, "参数错误"))
}
```

#### 日志
使用[https://github.com/uber-go/zap](https://github.com/uber-go/zap)实在日志组件，支持日志文件按日期自动分割。   
配置如下：
```yaml
server :
  logPath: /tmp/app.log #日志文件路径
  debug: true  # 是否开启debug日志
```

写日志
```go
import (
 "github.com/ebar-go/ego/component/log"
)
// 输出：{"level_name":"info","datetime":"2020-06-03 22:52:49","file":"log/log.go:24","message":"Info","system_name":"app","system_port":8080,"context":{"hello":"world","trace_id":""}}
log.Info("infoMessage", log.Context{
  "hello":"world",
})
log.Debug("debugMessage", log.Context{
  "hello":"world",
})
log.Error("errorMessage", log.Context{
  "hello":"world",
})
```

#### Mysql
基于`gorm`的数据库组件，支持多数据库连接，读写分离。

- 配置
```yaml
mysql:    # mysql配置。支持多数据库，读写分离
  default:  # 默认连接，如果还有其他数据库要连接，换个名字即可
    maxIdleConnections: 10  # 最大空闲连接数
    maxOpenConnections: 40  # 最大打开连接数
    maxLifeTime: 60          # 超时时间
    dsn:    # 连接配置，默认第一个为写库，也可以只配置一个，即读写使用一个连接
      - host: 127.0.0.1
        port: 32768
        user: root
        password: mysql
        name: blog
      - host: 127.0.0.1 # 从库,读库
        port: 32769
        user: root
        password: mysql
        name: blog
  otherDB:  # 其他数据库
    maxIdleConnections: 10  # 最大空闲连接数
    maxOpenConnections: 40  # 最大打开连接数
    maxLifeTime: 60          # 超时时间
    dsn:    # 连接配置，默认第一个为写库，也可以只配置一个，即读写使用一个连接
      - host: 127.0.0.1
        port: 32768
        user: root
        password: mysql
        name: blog
```

- 使用
```go
// 初始化数据库，一般在init函数里，加载了配置文件之后执行
secure.Panic(app.InitDB())
// 获取默认的库连接
app.DB()
// 指定库连接
app.GetDB("otherDB")
// 事务
if err := app.DB().Transaction(func(tx *gorm.Db) error {
    // some query or update
    return nil
}); err != nil {
    // 
    return fmt.Errorf("Save failed:%v", err)
}
```

### 事件
提供快捷的事件开发模式。
- 使用
```go
// 注册事件
event.Register("someEvent", event.Listener{
  Mode: event.Sync, // 同步事件,异步事件为 event.Async
  Handler: func(ev event.Event) {
    fmt.Println(ev.Params) 
  }
})
// 更快捷的注册事件
event.Listen("someEvent", func(ev event.Event) {
	fmt.Println(ev.Params)
})

// 触发事件
event.Trigger("someEvent", "someParam")
```
- 系统事件
```
// http服务启动前触发
event.BeforeHttpStart
// http服务启动后触发
event.AfterHttpStart
// http服务关闭前触发，平滑重启
event.BeforeHttpShutdown
// 数据库连接成功后触发
event.AfterDatabaseConnect
// 路由执行前触发
event.BeforeRoute
// 路由执行后触发
event.AfterRoute
```

### 集成Swagger，自动生成API文档。
刚开始我不喜欢让项目代码冗余太多得注释，但实际使用后，发现确实特别实用且方便。
关于swagger请参考我得文档，有详细说明，这里不再做过多描述。直接：[swagger使用说明](https://hongker.github.io/2020/06/10/golang-swagger/)