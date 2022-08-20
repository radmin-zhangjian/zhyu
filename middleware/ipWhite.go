package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhyu/setting"
)

func IpAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		flag := false
		ipList := setting.WhiteList.Ip
		for _, host := range ipList {
			if clientIP == host {
				flag = true
				break
			}
		}
		if flag == false {
			c.String(http.StatusUnauthorized, "ip %s not auth", clientIP)
			c.Abort()
		}
		c.Next()
	}
}
