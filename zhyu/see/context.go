package see

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"
)

const abortIndex int8 = math.MaxInt8 / 2

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	handlers   HandlersChain
	index      int8
	StatusCode int

	engine *Engine

	Keys map[string]any

	mu sync.RWMutex
}

func (c *Context) reset() {
	c.handlers = nil
	c.index = -1

	c.Keys = nil
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}

// PostForm 获取表单的值
func (c *Context) PostForm(key string) string {
	return c.Request.PostForm.Get(key)
}

// Query 获取url的值
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// setStatus设置状态码
func (c *Context) setStatus(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置头信息
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// String设置回复体
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-type", "text/plain")
	c.setStatus(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// Json 设置回复体
func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("Content-type", "application/json")
	c.setStatus(code)
	en := json.NewEncoder(c.Writer)
	if err := en.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// Html 设置回复体
func (c *Context) Html(code int, html string) {
	c.SetHeader("Content-type", "text/html")
	c.setStatus(code)
	c.Writer.Write([]byte(html))
}

// Data 设置回复体
func (c *Context) Data(code int, data []byte) {
	c.setStatus(code)
	c.Writer.Write(data)
}
