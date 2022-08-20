package main

import (
	"log"
	"net/http"
	"time"
	"zhyu/zhyu/see"
)

/*
 * 自定义简易脚手架
 * 自定义中间件
 */
func main() {
	// 初始化引擎
	r := see.New()
	// 中间件
	r.Use(Middle1())
	// 注册路由
	r.GET("/red", func(c *see.Context) {
		//c.Writer.Write([]byte("hello b ge"))
		data := map[string]any{
			"red": "我是red",
		}
		log.Printf("Red: %v", data)
		c.Json(http.StatusOK, data)
	})
	r.Use(Middle2())
	r.GET("/hello", Hello)
	r.GET("/say", Say)

	// 初始化
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8099",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func Hello(c *see.Context) {
	//fmt.Fprintf(c.Writer, "hello red")
	data := map[string]any{
		"hello": "我是Hello",
	}
	log.Printf("Hello: %v", data)
	c.Json(http.StatusOK, data)
}

func Say(c *see.Context) {
	data := map[string]any{
		"say": "我是say",
	}
	log.Printf("Say: %v", data)
	c.Json(http.StatusOK, data)
}

func Middle1() see.HandlerFunc {
	return func(c *see.Context) {
		log.Printf("Mid: %s", "我是中间件 Middle1 === start")
		c.Next()
		log.Printf("Mid: %s", "我是中间件 Middle1 === end")
	}
}

func Middle2() see.HandlerFunc {
	return func(c *see.Context) {
		log.Printf("Mid: %s", "我是中间件 Middle2 === start")
		c.Next()
		log.Printf("Mid: %s", "我是中间件 Middle2 === end")
	}
}
