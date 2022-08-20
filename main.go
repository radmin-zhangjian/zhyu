package main

import (
	"zhyu/servers"
)

func main() {
	http := servers.NewHttp()
	// 启动gin服务
	router := http.GinNew()
	// 启动http服务
	http.HttpServer(router)
}
