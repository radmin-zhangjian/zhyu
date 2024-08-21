package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"reflect"
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

// RateLimiter 限流
func RateLimiter(keys []string, value ...any) (int, error) {
	rdb := utils.GetRedis()

	// Lua 脚本
	luaScript := `
		local key = KEYS[1]
		local time = tonumber(ARGV[1])
		local count = tonumber(ARGV[2])
		local current = redis.call('get', key)
		if current and tonumber(current) > count then
		  return tonumber(current)
		end
		current = redis.call('incr', key)
		if tonumber(current) == 1 then
		  redis.call('expire', key, time)
		end
		return tonumber(current)
	`

	// 运行 Lua 脚本
	result, err := rdb.Eval(ctx, luaScript, keys, value).Result()
	if err != nil {
		fmt.Println("Error running Lua script:", err)
		return 0, errors.New(err.Error())
	}

	fmt.Println("Result from Lua script:", result)

	var num int
	switch reflect.TypeOf(result).Kind() {
	case reflect.Int:
		if i, ok := result.(int); ok {
			num = i
		}
	case reflect.Int64:
		if i, ok := result.(int64); ok {
			num = int(i)
		}
	case reflect.Float64:
		// 将 float64 转换为 int
		if f, ok := result.(float64); ok {
			i := int(f)
			num = i
		}
	}

	fmt.Println("Result from Lua num:", num)

	return num, nil
}
