package service

import (
	"zhyu/app/common"
	"zhyu/app/dao"
)

// RateLimiter 接口限流
func RateLimiter(keys []string, time int, count int) any {
	data := make(map[string]interface{})
	code := common.SUCCESS

	num, err := dao.RateLimiter(keys, time, count)
	if err != nil {
		code = common.ERROR
		return common.Result(code, common.GetMsg(code), data)
	}

	if num > count {
		code = common.RATE_LIMITER
	}

	return common.Result(code, common.GetMsg(code), data)
}
