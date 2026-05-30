package data

import (
	"gateway-service/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewUserRepo, NewGameAppRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	redis *redis.Client
}

func (d *Data) Redis() *redis.Client {
	return d.redis
}

// NewData .
func NewData(c *conf.Data, redisConfig *conf.RedisConfig) (*Data, func(), error) {
	redisAddr := "127.0.0.1:6379"
	readTimeout := 200 * time.Millisecond
	writeTimeout := 200 * time.Millisecond
	if c.Redis != nil {
		if c.Redis.GetAddr() != "" {
			redisAddr = c.Redis.GetAddr()
		}
		if c.Redis.GetReadTimeout() != nil {
			readTimeout = c.Redis.GetReadTimeout().AsDuration()
		}
		if c.Redis.GetWriteTimeout() != nil {
			writeTimeout = c.Redis.GetWriteTimeout().AsDuration()
		}
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		DB:           redisConfig.Database(),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})
	cleanup := func() {
		log.Info("closing the data resources")
		if err := redisClient.Close(); err != nil {
			log.Errorf("close redis client failed: %v", err)
		}
	}
	return &Data{redis: redisClient}, cleanup, nil
}
