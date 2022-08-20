package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/petermattis/goid"
	"log"
	"sync"
	"zhyu/app/common"
	"zhyu/utils"
	"zhyu/utils/logger"
)

// ContentKeys 定义全局的Keys中间件
func ContentKeys() gin.HandlerFunc {
	return func(c *gin.Context) {
		// uuid
		requestId := utils.Uuid.GetId()

		// 链路id
		c.Set("requestId", requestId)

		// context 上下文
		//ctx, cancel := context.WithCancel(c.Request.Context())
		ctx, _ := context.WithCancel(c.Request.Context())
		valueCtx := context.WithValue(ctx, "requestId", requestId)
		c.Set("ctx", valueCtx)
		//cancel()

		// 获取 goroutine id -- 为gorm的logger用
		goId := goid.Get()
		c.Set("goId", goId)

		var mu sync.RWMutex
		mu.Lock()
		common.RequestIdMap[goId] = requestId
		mu.Unlock()

		// 注册自定义logger
		logs := logger.Logger{Context: c}
		c.Set("logs", logs)

		c.Next()

		// 清空 RequestIdMap
		mu.Lock()
		delete(common.RequestIdMap, goId)
		mu.Unlock()
		log.Printf("common.RequestIdMap: %#v", common.RequestIdMap)
	}
}
