package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"zhyu/utils/logger"
)

type Context struct {
	*gin.Context
	Ctx      context.Context
	Logs     *logger.Logger
	UserInfo map[string]any
}

func NewApp() *Context {
	return new(Context)
}
