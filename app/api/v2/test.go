package v2

import (
	"log"
	"net/http"
	"time"
	"zhyu/app/service"
)

// Test 动态载入的api接口
// http://localhost:9090/api/v2/test
// 静态路由载入
// http://localhost:9090/v2/test
func (c *App) Test() {
	c.String(http.StatusOK, "Hello, v2 test")
}

func (c *App) SayDo() {
	time.Sleep(3 * time.Second)
	a := c.Query("a")
	b := c.PostForm("b")
	userId := c.UserInfo["Name"]
	log.Printf("UserInfo: %v", c.UserInfo)

	c.String(http.StatusOK, "Hello, v2 c.SayDo; %s, b: %s, userId: %s", a, b, userId)
}

func (c *App) SayIn() {
	log.Printf("UserInfo: %v", c.UserInfo)
	// GORM 测试
	resultDB := service.SayDbService(c.Context)
	c.JSON(http.StatusOK, resultDB)
}
