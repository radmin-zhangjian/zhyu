package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"time"
	"zhyu/app/common"
)

// 定义一个限流器，生成速率是800 个/s，令牌桶的容量是800。也就是每秒最多能通过800个请求。
var limiter = rate.NewLimiter(rate.Every(1*time.Second), 800)

// NewRateLimiter 限流
func NewRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		//limiter.SetLimit(1) // 动态设置
		//uri := fmt.Sprintf("%s://%s%s", c.Request.URL.Scheme, c.Request.Host, c.Request.URL.Path)
		if limiter.Allow() == false {
			//fmt.Println("服务繁忙，请稍后再试...", uri)
			mapData := common.Result(200, "服务繁忙，请稍后再试...", nil)
			c.JSON(200, mapData)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
