package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/url"
	"time"
	"zhyu/setting"
	"zhyu/utils"
	"zhyu/utils/logger"
)

// 全部级别
var logLevelArrayAll = []interface{}{
	"debug",
	"info",
	"warn",
	"error",
}

// info级别以上
var logLevelArrayInfo = []interface{}{
	"info",
	"warn",
	"error",
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (w responseBodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 response 内容
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// 获取请求数据
		var requestBody []byte
		if c.Request.Body != nil {
			// c.Request.Body 是一个 buffer 对象，只能读取一次
			requestBody, _ = ioutil.ReadAll(c.Request.Body)
			// 读取后，重新赋值 c.Request.Body ，以供后续的其他操作
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}

		// url.QueryEscape(urlStr)
		requestBodyStr, _ := url.QueryUnescape(string(requestBody))
		requestURI, _ := url.QueryUnescape(c.Request.RequestURI)

		startTime := time.Now()
		startTimeStr := startTime.Format("2006-01-02 15:04:05")
		// 赋值 可以用 c.MustGet("time") 获取
		c.Set("startTime", startTimeStr)

		c.Next()

		// 日志级别
		if utils.IsArray(setting.Server.LogLevel, logLevelArrayAll) {
			// 获取运行时间差
			//latency := time.Since(startTime).Microseconds() // 执行时间 微妙
			latency := float64(time.Since(startTime)) / 1e6 // 执行时间 毫秒
			//latency, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", latency), 64)
			latencyCost := fmt.Sprintf("%.3f", latency)
			//endTime := time.Now() // 结束时间
			//latency := endTime.Sub(startTime).Microseconds() // 执行时间
			//endTimeStr := endTime.Format("2006-01-02 15:04:05")

			// 响应数据
			responseBody := w.body.String()
			// 获取状态
			resposeStatus := c.Writer.Status()
			var logLevel string
			if resposeStatus > 400 && resposeStatus <= 499 {
				// 除了 StatusBadRequest 以外，warning 提示一下，常见的有 403 404，开发时都要注意
				logLevel = "WARNING"
			} else if resposeStatus >= 500 && resposeStatus <= 599 {
				// 除了内部错误，记录 error
				logLevel = "ERROR"
			} else {
				logLevel = "INFO"
			}

			logMsg := fmt.Sprintf("[%s][%s][%s][traceId:%v][host:%s][ip:%s][code:%d][cost:%s][%s %s %s %s][User-Agent:\"%s\"][request]%s[respose]%s[error]%s\n",
				setting.Server.ServerName,
				startTimeStr,
				logLevel,
				c.Keys["requestId"],
				c.Request.Host,
				c.ClientIP(),
				resposeStatus,
				latencyCost,
				c.Request.Method,
				//c.Request.URL.Path,
				requestURI,
				c.Request.Proto,
				c.Request.Header.Get("Content-Type"),
				c.Request.UserAgent(),
				requestBodyStr,
				responseBody,
				c.Errors.ByType(gin.ErrorTypePrivate).String(),
			)

			logs := c.MustGet("logs").(logger.Logger)
			logs.Print(logMsg)
		}
	}
}
