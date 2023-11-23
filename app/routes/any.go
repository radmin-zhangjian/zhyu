package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	app "zhyu/app"
	"zhyu/app/api"
	"zhyu/middleware"
	"zhyu/utils/logger"
)

// SortDesc 倒叙排序
func SortDesc(source []api.AppVersion) {
	sort.Slice(source, func(i, j int) bool {
		return source[i].Version > source[j].Version
	})
}

// InitRoutes 匹配路由
func InitRoutes(version string, controller string, c *gin.Context) interface{} {
	// 版本
	appData := api.Version()
	// 倒叙
	SortDesc(appData)

	// 强制版本分离
	//for k, v := range appData {
	//	log.Printf("k = %v, v = %v", k, v)
	//	if version == v.Version {
	//		return v.Object
	//	}
	//}

	// 向上版本查找
	var vi int
	for _, v := range appData {
		//log.Printf("k = %v, v = %v", k, v)
		ver1, _ := strconv.ParseInt(version[1:], 10, 64)
		ver2, _ := strconv.ParseInt(v.Version[1:], 10, 64)
		vi++
		if vi == 1 {
			//log.Printf("kk = %v, vv = %v", ver1, ver2)
			if ver1 > ver2 {
				return nil
			}
		}
		if ver1 >= ver2 {
			srv := v.Object
			rValue := reflect.ValueOf(srv)
			rType := reflect.TypeOf(srv)
			reciver := rValue.Elem().FieldByName("Context")
			reciver.Set(reflect.ValueOf(&app.Context{Context: c}))
			_, exist := rType.MethodByName(controller)
			if exist {
				return v.Object
			}
		}
	}

	for k, v := range appData {
		log.Printf("k = %v, v = %v", k, v)
		if version == v.Version {
			return v.Object
		}
	}

	return nil
}

// NewAny 路由动态载入 以api开始
// 例如：http://localhost:9090/api/v1/sayOut?a=abc
// 例如：http://localhost:9090/api/v2/sayOut?a=abc
func NewAny(router *gin.Engine) {
	anyRouter := router.Group("/api")
	anyRouter.Use(middleware.JwtAuth())
	{
		anyRouter.Any("/*action", func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, "/public/") {
				//匹配静态文件服务
			} else {
				path := strings.Split(c.Request.URL.Path, "/")
				version := strings.ToLower(path[2][:1]) + path[2][1:]
				controllerName := strings.ToUpper(path[3][:1]) + path[3][1:]

				srv := InitRoutes(version, controllerName, c)
				if srv == nil {
					c.String(http.StatusNotFound, "method %s not found!", c.Request.URL.Path)
					c.Abort()
					return
				}

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
				log.Printf("app pool: %v", &appc)
				appc.Reset()
				appc.Context = c
				appc.UserInfo = userInfo
				appc.Ctx = ctx
				appc.Logs = &logs
				defer appPool.Put(appc)
				//log.Printf("app pool222 ======: %v", &appc)

				rValue := reflect.ValueOf(srv)
				rType := reflect.TypeOf(srv)
				reciver := rValue.Elem().FieldByName("Context")
				// 原始模式
				//reciver.Set(reflect.ValueOf(&app.Context{Context: c, UserInfo: userInfo, Ctx: ctx, Logs: &logs}))
				// app new 连接池模式
				reciver.Set(reflect.ValueOf(appc))

				method, exist := rType.MethodByName(controllerName)
				if exist {
					args := []reflect.Value{rValue}
					method.Func.Call(args)
				} else {
					c.String(http.StatusNotFound, "method %s not found", c.Request.URL.Path)
				}
			}
		})
	}
}
