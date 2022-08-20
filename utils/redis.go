package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
	"zhyu/setting"
)

var rdb *redis.Client

func GetRedis() *redis.Client {
	return rdb
}

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     setting.Redis.Host + ":" + setting.Redis.Port,
		Password: setting.Redis.Password,
		DB:       int(setting.Redis.DB),
		PoolSize: int(setting.Redis.PoolSize), // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("redis connect get failed.%v", err)
		return
	}
}
