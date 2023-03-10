package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
	"zhyu/app/common"
	"zhyu/utils"
)

// 文档 https://github.com/go-redis/redis

var ctx = context.Background()

// SetString 存值
func SetString(key string, value any, exps ...int64) error {
	if key == "" || value == nil {
		return errors.New(common.GetMsg(common.INVALID_PARAMS))
	}
	var expiration time.Duration
	expiration = 0
	for exp := range exps {
		expiration = time.Duration(exp) // 过期时间
		break
	}
	rdb := utils.GetRedis()
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return errors.New(err.Error())
	}
	log.Printf("push key: %v, value: %v", key, value)
	return nil
}

// GetString 取值
func GetString(key string) (any, error) {
	if key == "" {
		return "", errors.New(common.GetMsg(common.INVALID_PARAMS))
	}
	rdb := utils.GetRedis()
	res, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New(fmt.Sprintf(common.GetMsg(common.INVALID_RESULT), key))
	}
	if err != nil {
		return "", errors.New(err.Error())
	}
	return res, nil
}
