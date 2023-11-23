package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"zhyu/app/common"
	"zhyu/app/service/auth"
	"zhyu/utils"
)

// JwtAuth JWT验证
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不需要验证的路由
		// todo
		if _, ok := common.RouteAuth[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		// 开始验证token
		var code int
		var data = make(map[string]interface{})
		var userId int64

		code = common.SUCCESS
		token := c.PostForm("token")
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			code = common.ERROR_AUTH_NO_TOKRN
		} else {
			claims, err := utils.ParseToken(token)
			log.Printf("claims: %#v", claims)
			if err == nil {
				userId = claims.UserId
			} else {
				code = common.ERROR_AUTH_CHECK_TOKEN_FAIL
				if claims != nil && time.Now().Unix() > claims.ExpiresAt {
					code = common.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				}
			}
		}

		// 如果token验证不通过，直接终止程序，c.Abort()
		if code != common.SUCCESS {
			data["code"] = code
			data["token"] = token
			// 返回错误信息
			c.JSON(http.StatusUnauthorized, data)
			//终止程序
			c.Abort()
			return
		}

		// 保存 userId
		c.Set("userId", userId)

		// 可选项 - 验证用户是否有效 根据情况 验证没问题了 也可以直接把用户信息放在keys里
		userInfo := auth.UserAuth(c.Request.Context(), userId)
		//log.Printf("jwt userInfo: %v", userInfo)
		if userInfo["code"] == common.SUCCESS {
			result := userInfo["data"].(map[string]map[string]any)
			c.Set("userName", result["userInfo"]["userName"])
		}

		c.Next()
	}
}
