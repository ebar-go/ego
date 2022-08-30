package example

import (
	"fmt"
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/component"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"net/http"
	"testing"
	"time"
)

func TestAggregatorWithComponents(t *testing.T) {
	aggregator := ego.NewAggregatorServer()

	// component logger
	component.Provider().Logger().Info("test logger info function")

	// component cache
	component.Provider().Cache().Default().Set("someCacheKey", "someCacheValue", time.Minute)

	// component redis
	if err := component.Provider().Redis().Connect(&redis.Options{
		// ... some options like address, port
	}); err != nil {
		panic(err)
	}
	component.Provider().Redis().Set("someRedisKey", "someRedisVal", time.Minute)

	// 实例化一个http服务
	httpServer := ego.NewHTTPServer(":8080").
		RegisterRouteLoader(func(router *gin.Engine) { // 注册路由
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "home")
			})
		})

	aggregator.WithServer(httpServer)

	aggregator.Run()
}

type User struct {
	Id   int
	Name string
}

func (u *User) TableName() string {
	return "users"
}

type OtherUser struct {
	Id   int
	Name string
}

func (u *OtherUser) TableName() string {
	return "other_users"
}
func TestComponentGorm(t *testing.T) {
	// test primary connection
	dsn := "root:123456@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
	err := component.Provider().Gorm().OpenMySQL(dsn)
	assert.Nil(t, err)

	// test insert
	insertErr := component.Provider().Gorm().Create(&User{Name: "John"}).Error
	assert.Nil(t, insertErr)

	// use other connection
	otherDsn := "root:123456@tcp(127.0.0.1:3306)/db2?charset=utf8mb4&parseTime=True&loc=Local"
	resolverErr := component.Provider().Gorm().RegisterResolverConfig(dbresolver.Config{
		Sources: []gorm.Dialector{mysql.Open(otherDsn)},
	}, "other_users")
	assert.Nil(t, resolverErr)

	// insert other users with db2 connection
	otherInsertErr := component.Provider().Gorm().Create(&OtherUser{Name: "Tom"}).Error
	assert.Nil(t, otherInsertErr)

	// test insert primary connection
	insertPrimaryErr := component.Provider().Gorm().Create(&User{Name: "John"}).Error
	assert.Nil(t, insertPrimaryErr)
}

func TestComponentConfig(t *testing.T) {
	conf := component.Provider().Config()
	err := conf.LoadFile("example.yaml")
	assert.Nil(t, err)

	fmt.Println(conf.GetInt("app.someInt"))
	fmt.Println(conf.GetBool("app.someBool"))
	fmt.Println(conf.GetString("app.someString"))
}

func TestComponentCurl(t *testing.T) {
	curl := component.Provider().Curl()

	response, err := curl.Get("http://localhost:8080/")
	assert.Nil(t, err)
	fmt.Println(string(response.Bytes()))
	response.Release()
}

func TestComponentRedis(t *testing.T) {
	// connect redis
	if err := component.Provider().Redis().Connect(&redis.Options{Addr: "127.0.0.1:6379"}); err != nil {
		t.Fatal(err)
	}

	// set some key
	component.Provider().Redis().Set("someKey", "someValue", time.Second)

}
