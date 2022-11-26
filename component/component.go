package component

import (
	"github.com/ebar-go/ego/component/cache"
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/curl"
	"github.com/ebar-go/ego/component/db"
	"github.com/ebar-go/ego/component/event"
	"github.com/ebar-go/ego/component/jwt"
	"github.com/ebar-go/ego/component/logger"
	"github.com/ebar-go/ego/component/mongo"
	"github.com/ebar-go/ego/component/redis"
	"github.com/ebar-go/ego/component/tracer"
	"github.com/ebar-go/ego/component/validator"
	"github.com/ebar-go/ego/utils/structure"
	gocache "github.com/patrickmn/go-cache"
)

var (
	cacheBuilder      = structure.NewSingleton[*cache.Builder](cache.New)
	configInstance    = structure.NewSingleton[*config.Instance](config.New)
	curlInstance      = structure.NewSingleton[*curl.Instance](curl.New)
	eventInstance     = structure.NewSingleton[*event.Instance](event.New)
	jwtInstance       = structure.NewSingleton[*jwt.Instance](jwt.New)
	tracerInstance    = structure.NewSingleton[*tracer.Instance](tracer.New)
	validatorInstance = structure.NewSingleton[*validator.Instance](validator.New)
	loggerInstance    = structure.NewSingleton[*logger.Instance](logger.New)
	dbInstance        = structure.NewSingleton[*db.Instance](db.New)
	redisInstance     = structure.NewSingleton[*redis.Instance](redis.New)
	mgo               = structure.NewSingleton[*mongo.Instance](mongo.New)
)

func Cache() *gocache.Cache {
	return cacheBuilder.Get().Default()
}

func CacheBuilder() *cache.Builder {
	return cacheBuilder.Get()
}

func Config() *config.Instance {
	return configInstance.Get()
}

func Curl() *curl.Instance {
	return curlInstance.Get()
}

func Event() *event.Instance {
	return eventInstance.Get()
}

func JWT() *jwt.Instance {
	return jwtInstance.Get()
}

func Tracer() *tracer.Instance {
	return tracerInstance.Get()
}

func Validator() *validator.Instance {
	return validatorInstance.Get()
}

func Logger() *logger.Instance {
	return loggerInstance.Get()
}

func DB() *db.Instance {
	return dbInstance.Get()
}

func Redis() *redis.Instance {
	return redisInstance.Get()
}

func Mgo() *mongo.Instance {
	return mgo.Get()
}
