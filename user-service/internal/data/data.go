package data

import (
	"time"
	"user-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet 注册数据层依赖，供 Wire 自动注入使用。
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewUserRepo)

// Data 保存数据层共享资源，例如数据库连接池。
type Data struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewData 初始化数据层资源，创建 GORM 数据库连接并返回清理函数。
func NewData(c *conf.Data, redisConfig *conf.RedisConfig) (*Data, func(), error) {
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	redisDB := 0
	if redisConfig != nil {
		redisDB = redisConfig.DB
	}
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

	// 创建 Redis 客户端，使用配置中的地址、数据库索引和超时设置。
	redisClient := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		DB:           redisDB,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})
	cleanup := func() {
		log.Info("closing the data resources")
		sqlDB, err := db.DB()
		if err != nil {
			log.Errorf("get database connection failed: %v", err)
			return
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
		if err := sqlDB.Close(); err != nil {
			log.Errorf("close database connection failed: %v", err)
		}
		if err := redisClient.Close(); err != nil {
			log.Errorf("close redis client failed: %v", err)
		}
	}
	return &Data{db: db, redis: redisClient}, cleanup, nil
}
