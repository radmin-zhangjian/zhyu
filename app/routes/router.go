package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"sync"
	app "zhyu/app"
	"zhyu/app/api/auth"
	v1 "zhyu/app/api/v1"
	v2 "zhyu/app/api/v2"
	"zhyu/middleware"
	"zhyu/utils/logger"
)

var appPool = sync.Pool{
	New: func() any {
		return app.NewApp()
	},
}

type HandlerFunc func(ctx *app.Context)

// 函数方法调用 通过HandlerFunc实现
// 例如：V1.GET("/sayIn", handler(v1.SayIn))
func handler(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 可以把 keys 下的数据统一放在 app.context 方便调用
		userInfo := make(map[string]any)
		if userId, ok := c.Keys["userId"]; ok {
			userInfo["id"] = userId.(int64)
		}
		if userName, ok := c.Keys["userName"]; ok {
			userInfo["userName"] = userName.(string)
		}
		// context 上下文
		ctx := c.MustGet("ctx").(context.Context)
		// 自定义 logs
		logs := c.MustGet("logs").(logger.Logger)

		content := app.NewApp()
		content.Context = c
		content.UserInfo = userInfo
		content.Ctx = ctx
		content.Logs = &logs
		handler(content)
	}
}

// 对象方法调用 通过反射实现
// 例如：V1.GET("/sayIn", handlerRef("SayIn", v1.New()))
func handlerRef(api string, object any) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 可以把 keys 下的数据统一放在 app.context 方便调用
		userInfo := make(map[string]any)
		if userId, ok := c.Keys["userId"]; ok {
			userInfo["id"] = userId.(int64)
		}
		if userName, ok := c.Keys["userName"]; ok {
			userInfo["userName"] = userName.(string)
		}
		// context 上下文
		ctx := c.MustGet("ctx").(context.Context)
		// 自定义 logs
		logs := c.MustGet("logs").(logger.Logger)

		// app new 连接池模式
		appc := appPool.Get().(*app.Context)
		//log.Printf("app pool: %v", appc)
		appc.Reset()
		appc.Context = c
		appc.UserInfo = userInfo
		appc.Ctx = ctx
		appc.Logs = &logs
		defer appPool.Put(appc)
		//log.Printf("app pool222 ======: %v", &appc)

		srv := object
		rValue := reflect.ValueOf(srv)
		rType := reflect.TypeOf(srv)
		reciver := rValue.Elem().FieldByName("Context")
		// 原始模式
		//reciver.Set(reflect.ValueOf(&app.Context{Context: c, UserInfo: userInfo, Ctx: ctx, Logs: &logs}))
		// app new 连接池模式
		reciver.Set(reflect.ValueOf(appc))
		method, exist := rType.MethodByName(api)
		if exist {
			args := []reflect.Value{rValue}
			method.Func.Call(args)
		} else {
			c.String(http.StatusNotFound, "method %s not found", c.Request.URL.Path)
		}
	}
}

// Routes 静态路由
func Routes(router *gin.Engine) {

	//router.POST("/login", handler(auth.Login))
	router.POST("/login", handlerRef("Login", auth.New()))
	router.GET("/middleware", handlerRef("Middleware", auth.New()))

	// 测试
	router.GET("/v1/estest", handlerRef("EsTest", v1.New()))
	router.GET("/v1/test", handlerRef("Test", v1.New()))
	router.GET("/v2/test", handlerRef("Test", v2.New()))

	authorized := router.Group("/")
	authorized.Use(middleware.JwtAuth())
	{
		// v1
		V1 := authorized.Group("v1")
		{
			//V1.GET("/sayIn", handler(v1.SayIn))
			//V1.POST("/sayDo", handler(v1.SayDo))
			V1.Any("/sayIn", handlerRef("SayIn", v1.New()))
			V1.Any("/sayDo", handlerRef("SayDo", v1.New()))
			user := V1.Group("user")
			{
				// visit 0.0.0.0:8080/user/sayIn
				user.Any("/sayIn", handlerRef("SayIn", v1.New()))
			}
		}

		// v2
		V2 := authorized.Group("v2")
		{
			V2.Any("/sayIn", handlerRef("SayIn", v2.New()))
			V2.Any("/sayDo", handlerRef("SayDo", v2.New()))
		}
	}

}
