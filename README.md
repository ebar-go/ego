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

#### 加载配置
```go
// 从环境变量中读取
config.ReadFromEnvironment()
// 或者从文件读取
// config.ReadFromFile(configFilePath)

```
#### http请求
使用原生http
```go
package main
import (
    "fmt"
    "net/http"
    "github.com/ebar-go/ego/utils/curl"
)

func main() {
    address := "http://baidu.com"
    request,_ := http.NewRequest(http.MethodGet, address, nil)
    resp, err := curl.Execute(request)
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
    // RequestLog middleware
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
```go
package main
import (
 "github.com/ebar-go/ego/http/paginator"
)

func main() {
    pagination := paginator.Paginate(20, 1, 10)
    paginationSlice := paginator.PaginateSlice([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, 1, 10)
}
```

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
集成`github.com/go-playground/validator`,支持自定义字段名称`comment`

- 使用方式:
```go
package main
import (
"fmt"
"github.com/ebar-go/ego/http/validator"
"github.com/gin-gonic/gin"
"github.com/gin-gonic/gin/binding"
)
func main() {
    // 设置自定义验证器
    binding.Validator = new(validator.Validator)
    router := gin.Default()
    
    	router.POST("user/auth", func(ctx *gin.Context) {
    		type AuthRequest struct {
    			// 如果是表单提交，使用form,否则获取不到数据
    			Email string `json:"email" validate:"required,email" comment:"邮箱"` // 验证邮箱格式
    			Pass string `json:"pass" binding:"required,min=6,max=10" comment:"密码"` // 验证密码，长度为6~10
    		}
    
    		var request AuthRequest
    		// 使用bind
    		if err := ctx.ShouldBindJSON(&request); err != nil {
    			ctx.JSON(200, gin.H{
    				"message" : err.Error(),
    			})
    			return
    		}
    
    		// other logic..
    	})
    
    	router.Run(":8080")
}
```

### 组件包
#### apollo
```go
_ := apollo.Init(apollo.Conf{
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
基于`go.uber.org/zap`实现的按日期分割的日志组件

- 使用方式:
```go
package main
import "github.com/ebar-go/ego/component/log"
func main() {
    log.App().Info("infoContent", log.Context(map[string]interface{}{
        "hello": "world",
    }))
}
```
#### mns
阿里云mns客户端
```go
// trigger mns init event
// make sure config is not empty 
event.DefaultDispatcher().Trigger(app.MNSClientInitEvent, nil)
mns := app.Mns()
```

#### mysql
- connect database
```go
// trigger mysql connect event
// make sure config is not empty 
event.DefaultDispatcher().Trigger(app.MySqlConnectEvent, nil)
db := app.Mysql() 
// .. curd
```

#### redis
- connect
```go
event.DefaultDispatcher().Trigger(app.RedisConnectEvent, nil)
redis := app.Redis()
```
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
