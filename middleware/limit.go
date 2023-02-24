package middleware

import (
	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
)

func LimitHandler() gin.HandlerFunc {
	// 每秒链接次数
	lmt := tollbooth.NewLimiter(10, nil)
	lmt.SetMessage("服务繁忙，请稍后再试...")
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}
