package auth

import (
	"log"
	"math"
	"net/http"
)

const abortIndex int8 = math.MaxInt8 / 2

type HandlerFunc func(*context)
type HandlersChain []HandlerFunc

type context struct {
	Handlers HandlersChain
	Index    int8
	Keys     []map[string]any
}

func (m *context) Next() {
	m.Index++
	for m.Index < int8(len(m.Handlers)) {
		m.Handlers[m.Index](m)
		m.Index++
	}
}

func (m *context) Abort() {
	m.Index = abortIndex
}

// Middleware 中间件的测试
func (c *App) Middleware() {
	// see 中间件实例
	m := &context{Index: -1}

	// 注册 func
	m.Handlers = make(HandlersChain, 0)
	m.Handlers = append(m.Handlers, m1())
	m.Handlers = append(m.Handlers, m2())
	m.Handlers = append(m.Handlers, m3())
	m.Handlers = append(m.Handlers, action)

	// 限制handler个数
	if len(m.Handlers) >= int(abortIndex) {
		panic("too many handlers")
	}

	// 开始运行
	m.Next()

	// echo 中间件实例
	echo()

	c.JSON(http.StatusOK, "中间件测试")
}

func action(m *context) {
	log.Println("main handler")
}

func m1() HandlerFunc {
	return func(m *context) {
		log.Println("m1 start")
		m.Next()
	}
}

func m2() HandlerFunc {
	return func(m *context) {
		log.Println("m2 start")
		//m.Abort()
		m.Next()
		log.Println("m2 end")
	}
}

func m3() HandlerFunc {
	return func(m *context) {
		log.Println("m3 start")
		m.Next()
		log.Println("m3 end")
	}
}

// =====================
// echo middleware 实例
// 这个其实 高阶函数和闭包的妙用  延迟执行
// 可以把函数作为参数传递给另一个函数，这样可以在调用这个函数时才会执行它。
// 通过简单的中间件 实现了装饰器模式
// 类似套娃或者洋葱模型  重上往下
// 最后调用的时候 会先经过装饰器的函数  再到真正要执行的函数
// 装饰器函数的循环调用 有个很秒的地方   第一次循环的时候  把真正要要执行的函数包裹起来 变成返回的函数 第二次执行的时候 又包裹了一层  第n次执行的时候 就继续包裹 直到结束
// 执行的时候 就是类似剥洋葱一样 一层层执行
// =====================
type fn func(string)
type fn2 func(next fn) fn
type middleList struct {
	handlers []fn2
}

func echo() {
	e := &middleList{handlers: []fn2{}}
	e.handlers = append(e.handlers, midd1)
	e.handlers = append(e.handlers, midd2)

	h := middleTofn(e, h)
	h("hello world")
}

func middleTofn(middleList *middleList, s fn) fn {
	for i := len(middleList.handlers) - 1; i >= 0; i-- {
		s = middleList.handlers[i](s)
	}
	return s
}

func h(string2 string) {
	log.Println("------------------" + string2)
}

func midd1(next fn) fn {
	return func(cxt string) {
		log.Println("midd1 start")
		next(cxt)
		log.Println("midd1 end")
	}
}

func midd2(next fn) fn {
	return func(cxt string) {
		log.Println("midd2 start")
		next(cxt)
		log.Println("midd2 end")
	}
}
