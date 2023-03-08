package middleware

import (
	"encoding/json"
	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
	"zhyu/app/common"
)

// LimitHandler 限流
func LimitHandler() gin.HandlerFunc {
	// 每秒链接次数
	lmt := tollbooth.NewLimiter(800, nil)
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			mapData := common.Result(200, "服务繁忙，请稍后再试...", nil)
			bytes, _ := json.Marshal(mapData)
			stringData := string(bytes)
			lmt.SetMessage(stringData)
			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}
