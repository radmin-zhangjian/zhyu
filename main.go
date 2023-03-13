package main

import (
	"zhyu/servers"
)

// 带上环境变量打包
// CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o app_zhyu .

func main() {
	http := servers.NewHttp()
	// 启动gin服务
	router := http.GinNew()
	// 启动http服务
	http.HttpServer(router)
}
