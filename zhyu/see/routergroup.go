package see

import (
	"log"
	"sync"
)

// RouterGroup 路由结构体
type RouterGroup struct {
	HandlersMap map[string]HandlersChain
	Handlers    HandlersChain
	engine      *Engine
	mu          sync.RWMutex
}

// addRoute 框架新增路由
func (group *RouterGroup) addRoute(method string, pattern string, handlers HandlersChain) {
	group.mu.Lock()
	defer group.mu.Unlock()
	log.Printf("Route %s - %s", method, pattern)
	key := method + "-" + pattern
	group.HandlersMap[key] = handlers
}

// Use adds middleware to the group.
func (group *RouterGroup) Use(middleware ...HandlerFunc) {
	group.Handlers = append(group.Handlers, middleware...)
}

// handle 每次请求需要的 handler
func (group *RouterGroup) handle(method string, pattern string, handler HandlersChain) {
	finalSize := len(group.Handlers) + len(handler)
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handler)
	group.addRoute(method, pattern, mergedHandlers)
}

// GET 匹配get方法
func (group *RouterGroup) GET(pattern string, handler ...HandlerFunc) {
	group.handle("GET", pattern, handler)
}

// POST 匹配post方法
func (group *RouterGroup) POST(pattern string, handler ...HandlerFunc) {
	group.handle("POST", pattern, handler)
}
