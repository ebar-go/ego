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
    
    utils.Panic("StartServer", server.Start())
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

#### 分页组件
对数据进行分页

#### 响应
固定格式的json输出

#### validator
集成`github.com/go-playground/validator`

### 组件包
#### apollo
#### auth
#### consul
微服务(SOA),集成consul组件
- 服务注册
- 服务发现

#### 日志
#### mns
阿里云mns客户端

#### mysql
自定义model

#### prometheus
集成prometheus

#### trace
全局ID

### 配置项
读取配置项

### utils
工具库
- 数组
- 转换
- 日期
- 文件
- json
- 数字
- 字符串

### log
日志管理器
- 系统日志
