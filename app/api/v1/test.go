package v1

import (
	"context"
	"log"
	"net/http"
	"time"
	"zhyu/app"
	"zhyu/app/service"
	v1 "zhyu/app/service/v1"
	"zhyu/utils"
	"zhyu/utils/logger"
	"zhyu/utils/uuid"
)

// SayDo 静态接口通过 handler 实现的方法
func SayDo(c *app.Context) {
	time.Sleep(3 * time.Second)
	a := c.Query("a")
	b := c.PostForm("b")
	userId := c.UserInfo["Name"]
	log.Printf("UserInfo: %v", c.UserInfo)
	c.String(http.StatusOK, "Hello, v1 SayDo; %s, b: %s, userId: %s", a, b, userId)
}

// SayIn 静态接口通过 handler 实现的方法
func SayIn(c *app.Context) {
	a := c.Query("a")
	log.Printf("UserInfo: %v", c.UserInfo)
	c.String(http.StatusOK, "Hello, v1 SayIn; %s", a)
}

// Test 动态载入的api接口
// http://localhost:9090/api/v1/test
// 静态路由载入
// http://localhost:9090/v1/test
func (c *App) Test() {
	c.String(http.StatusOK, "Hello, v1 test")
}

// EsTest elastic测试
func (c *App) EsTest() {
	// 插入
	//result := v1.EsCreateService(c.Context)
	// 修改
	//result := v1.EsUpdateService(c.Context)
	// 查询单个
	//result := v1.EsGetService(c.Context)
	// 查询多个
	//result := v1.EsSearchService(c.Context)
	// query查询
	result := v1.EsQueryService(c.Context)

	c.JSON(http.StatusOK, result)
}

// SayIn 动态载入的api接口
// http://localhost:9090/api/v1/sayIn
func (c *App) SayIn() {
	// 获取中间件赋值的参数
	timeNow := c.MustGet("startTime").(string)

	// get参数
	a := c.Query("a")
	// post参数
	b := c.DefaultPostForm("b", "bbb")

	//ctx1, cancel := context.WithCancel(c.Ctx)
	//ctx1.Done()
	//cancel()

	log.Printf("c.ctx.requestId: %v", c.Ctx.Value("requestId"))
	ctx, _ := c.MustGet("ctx").(context.Context)
	log.Printf("get.ctx: %v", ctx.Value("requestId"))
	log.Printf("c.keys.requestId: %v", c.Keys["requestId"])
	log.Printf("c.keys.UserId: %v", c.Keys["UserId"])

	logger.Info("uuid0:%v", utils.Uuid.GetId())
	// 雪花算法 获取uuid
	c.Logs.Info("uuid1:%v", utils.Uuid.GetId())
	time.Sleep(3 * time.Second)
	c.Logs.Info("asddds")
	c.Logs.Info("uuid2:%v", utils.Uuid.GetId())
	c.Logs.Info("uuid3:%v", utils.Uuid.GetId())

	// 索尼雪花算法
	sonyflake := utils.Uuid.Sonyflake()
	c.Logs.Info("sony uuid1:%v", uuid.NextID(sonyflake))
	c.Logs.Info("sony uuid2:%v", uuid.NextID(sonyflake))

	c.String(http.StatusOK, "Hello, c.v1 SayIn; query: %s; time: %s; b: %s", a, timeNow, b)
}

// SayDo 动态载入的api接口
// http://localhost:9090/api/v1/sayOut
func (c *App) SayDo() {
	// redis 测试
	result := service.SayRedis(c.Context)

	// 区分版本的测试
	resultV1 := v1.Say(c.Context)

	// GORM 测试
	resultDB := service.SayDb(c.Context)

	c.String(http.StatusOK, "c.v1 result: %s, resultV1: %s, UserInfo: %s", result, resultV1, c.UserInfo)
	c.JSON(http.StatusOK, resultDB)
}
