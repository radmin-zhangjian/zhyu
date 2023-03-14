package servers

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"time"
	"zhyu/app/routes"
	"zhyu/middleware"
	"zhyu/setting"
	"zhyu/utils"
	"zhyu/utils/logger"
	//_ "zhyu/setting/toml"
	//_ "zhyu/setting/etcd"
)

type Http struct {
}

func NewHttp() *Http {
	return &Http{}
}

// GinNew 初始化gin
func (s *Http) GinNew() *gin.Engine {
	// 启动模式
	gin.SetMode(setting.Server.RunMode)

	if "debug" == setting.Server.RunMode {
		// 日志始终着色
		gin.ForceConsoleColor()
		// 将日志写入文件
		f, _ := os.Create(setting.Server.LogPath + "/see.log")
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)      // 日志信息
		gin.DefaultErrorWriter = io.MultiWriter(f, os.Stdout) // 错误信息
	}

	// 没有中间件的引擎
	router := gin.New()
	if "debug" == setting.Server.RunMode {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())

	// Content Keys 需要放在路由前面进行初始化
	router.Use(middleware.ContentKeys())
	// 自定义Logger
	router.Use(middleware.Logger())
	go logger.LogHandlerFunc() // 异步处理日志
	// ip白名单
	router.Use(middleware.IpAuth())
	// 定义全局的CORS中间件
	router.Use(middleware.Cors())

	// 静态资源加载，css,js以及资源图片
	//router.StaticFS("/public", http.Dir("./website/static"))
	//router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// 导入所有模板
	//router.LoadHTMLGlob("website/tpl/*")

	// 限流 tollbooth-limit 中间件
	//router.Use(middleware.LimitHandler())
	// 限流 rate-limit 中间件
	router.Use(middleware.NewRateLimiter())

	// 注册静态路由
	routes.Routes(router)

	// 注册动态路由 以api开头
	routes.NewAny(router)

	// redis 初始化
	utils.InitRedis()
	// gorm 初始化
	utils.InitDB()
	// es 初始化
	utils.InitES()
	// 开启etcd 并 监听降级服务开关
	go utils.NewEtcd().ListenReduceRank()

	//routes.Run(":9090") // listen and serve on 0.0.0.0:8080
	return router
}

// HttpServer 启动服务 & 优雅Shutdown（或重启）服务
func (s *Http) HttpServer(router *gin.Engine) {
	// endless 热启动
	els := endless.NewServer(":"+setting.Server.Port, router)
	els.ReadHeaderTimeout = time.Duration(setting.Server.ReadTimeout) * time.Second
	els.WriteTimeout = time.Duration(setting.Server.WriteTimeout) * time.Second
	//go func() {
	err := els.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
	log.Println("Server on " + setting.Server.Port + " stopped")
	//}()

	//srv := &http.Server{
	//	Addr:         ":" + setting.Server.Port,
	//	Handler:      router,
	//	ReadTimeout:  time.Duration(setting.Server.ReadTimeout) * time.Second,
	//	WriteTimeout: time.Duration(setting.Server.WriteTimeout) * time.Second,
	//}
	//go func() {
	//	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		log.Fatalf("listen: %s\n", err)
	//	}
	//}()

	// 5秒后优雅Shutdown服务
	//quit := make(chan os.Signal)
	//signal.Notify(quit, os.Interrupt) //syscall.SIGKILL
	//<-quit
	//log.Println("Shutdown Server ...", setting.Server.ShutdownTime)
	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(setting.Server.ShutdownTime)*time.Second)
	//defer cancel()
	//if err := srv.Shutdown(ctx); err != nil {
	//	log.Fatal("Server Shutdown:", err)
	//}
	//select {
	//case <-ctx.Done():
	//}
	//log.Println("Server exiting")
}
