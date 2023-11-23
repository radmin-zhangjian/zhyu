package auth

import (
	"net/http"
	"time"
	"zhyu/app/common"
)

// TokenAuth token验证
func (c *App) TokenAuth() {
	// auth
	data := make(map[string]interface{})
	code := common.SUCCESS
	dateTime := time.Now().Unix()
	//dateTime := time.Now().Format("2006-01-02 15:04:05")  // 格式化
	data["dateTime"] = dateTime
	data = common.Result(code, common.GetMsg(code), data)
	c.JSON(http.StatusOK, data)
}
