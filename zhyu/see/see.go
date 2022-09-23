package see

import (
	"net/http"
	"sync"
)

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc

type Engine struct {
	RouterGroup

	pool sync.Pool
}

func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			//basePath: "/",
			//root:     true,
			HandlersMap: make(map[string]HandlersChain),
		},
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() any {
		return engine.allocateContext()
	}

	return engine
}

func (engine *Engine) allocateContext() *Context {
	return &Context{engine: engine}
}

// ServeHTTP 实现ServeHTTP方法 根据请求的方法及路径来匹配Handler
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := engine.pool.Get().(*Context)
	c.Writer = w
	c.Request = r
	c.reset()

	engine.handleHTTPRequest(c)
	engine.pool.Put(c)
}

// HandleContext重新进入已重写的上下文。
// 这可以通过将c.Request.URL.Path设置为新目标来实现。
// 免责声明:你可以循环自己来处理这个问题，明智地使用。
func (engine *Engine) handleHTTPRequest(c *Context) {
	key := c.Request.Method + "-" + c.Request.URL.Path
	if handlers, ok := engine.HandlersMap[key]; ok {
		c.handlers = handlers
	} else {
		c.String(http.StatusNotFound, "404 not found %s ", c.Request.URL.Path)
	}

	c.Next()
}

// Run 启动服务
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
