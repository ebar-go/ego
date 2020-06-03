## 说明
基于Gin的web微服务框架。

## 特性
- 简单的HTTP服务启动方式。
- 提供强大的日志管理器，支持按日期自动分割。
- 丰富的中间件(请求日志、JWT认证,跨域,Recover,全局链路)
- 支持基于consul,etcd的服务注册、发现、注销
- 全局Mysql,Redis连接池
- 自定义参数验证器
- 配置文件,支持json,ini,yaml等格式

## 安装

```
go get -u github.com/ebar-go/ego
```

## 示例
完整的web项目请查看[github.com/ebar-go/ego-demo](github.com/ebar-go/ego-demo)

## 模块
### web

- 启动web服务

```go
package main
import (
    "fmt"
    "github.com/ebar-go/ego/http"
    "github.com/gin-gonic/gin"
    "github.com/ebar-go/ego/utils/secure"
)
func main() {
    server := http.NewServer()
    // 添加路由
    server.Router.GET("/test", func(context *gin.Context) {
        fmt.Println("hello,world")
    })
    secure.Panic(server.Start())
}
```

#### 加载配置
```go
// 从环境变量中读取
config.ReadFromEnvironment()
// 或者从文件读取
config.ReadFromFile(configFilePath)

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

#### validator
集成`github.com/go-playground/validator`,支持自定义字段名称`comment`

- 使用方式:
```go
// 给Gin设置全局自定义验证器
binding.Validator = new(validator.Validator)
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
