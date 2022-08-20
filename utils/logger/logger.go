package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/petermattis/goid"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"zhyu/app/common"
	"zhyu/setting"
)

type Logger struct {
	*gin.Context
}

var logChannel = make(chan map[string]interface{}, 300)

// CreateDateDir 根据时间检测目录，不存在则创建
func CreateDateDir(Path string) string {
	folderName := time.Now().Format("20060102")
	folderPath := filepath.Join(Path, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.MkdirAll(folderPath, os.ModePerm)
		os.Chmod(folderPath, os.ModePerm)
	}
	return folderPath
}

// FilePath 文件名
func FilePath() string {
	rootPath, _ := os.Getwd() //获取项目根路径
	nowPath := rootPath + "/" + setting.Server.LogPath
	//pathArr := strings.Split(thisPath, "/")
	//fileName := strings.Join(pathArr, "-")
	//fileName = path.Join(writePath, fileName[1:len(fileName)]+".log")
	writePath := CreateDateDir(nowPath) // 根据时间检测是否存在目录，不存在创建
	fileName := setting.Server.ServerName + "_" + time.Now().Format("20060102") + ".log"
	fileName = path.Join(writePath, fileName)
	return fileName
}

// WriteWithIo 使用io.WriteString()函数进行数据的写入，不存在则创建
func WriteWithIo(content string) {
	filePath := FilePath()
	fileObj, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("Failed to open the file", err.Error())
		os.Exit(2)
	}
	if _, err := io.WriteString(fileObj, content); err == nil {
		fmt.Println("Successful appending to the file with os.OpenFile and io.WriteString.\n", content)
	}
}

// LogHandlerFunc 异步处理日志
func LogHandlerFunc() {
	for logField := range logChannel {
		var msgStr string
		if msg, ok := logField["msg"]; ok {
			msgStr = msg.(string)
		}
		WriteWithIo(msgStr)
	}
}

func join(logLevel string, message interface{}) string {
	startTimeStr := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[%s][%s][%s][traceId:%v][respose]%s\n",
		setting.Server.ServerName,
		startTimeStr,
		logLevel,
		common.RequestIdMap[goid.Get()],
		message,
	)
	return logMsg
}

// Debug 调试日志
func Debug(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	Print(join("DEBUG", message))
}

// Info 告知类日志
func Info(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	Print(join("INFO", message))
}

// Warn 警告类
func Warn(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	Print(join("WARNING", message))
}

// Error 错误时记录，不应该中断程序，查看日志时重点关注
func Error(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	Print(join("ERROR", message))
}

// Fatal 级别同 Error(), 写完 log 后调用 os.Exit(1) 退出程序
func Fatal(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	Print(join("FATAL", message))
}

func (log *Logger) join(logLevel string, message interface{}) string {
	startTimeStr := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[%s][%s][%s][traceId:%v][host:%s][ip:%s][code:%d][%s %s %s %s][respose]%s[error]%s\n",
		setting.Server.ServerName,
		startTimeStr,
		logLevel,
		log.Keys["requestId"],
		log.Request.Host,
		log.ClientIP(),
		log.Context.Writer.Status(),
		log.Request.Method,
		//log.Request.URL.Path,
		log.Request.RequestURI,
		log.Request.Proto,
		log.Request.Header.Get("Content-Type"),
		message,
		log.Errors.ByType(gin.ErrorTypePrivate).String(),
	)
	return logMsg
}

// Debug 调试日志
func (log *Logger) Debug(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	log.Print(log.join("DEBUG", message))
}

// Info 告知类日志
func (log *Logger) Info(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	log.Print(log.join("INFO", message))
}

// Warn 警告类
func (log *Logger) Warn(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	log.Print(log.join("WARNING", message))
}

// Error 错误时记录，不应该中断程序，查看日志时重点关注
func (log *Logger) Error(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	log.Print(log.join("ERROR", message))
}

// Fatal 级别同 Error(), 写完 log 后调用 os.Exit(1) 退出程序
func (log *Logger) Fatal(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	log.Print(log.join("FATAL", message))
}

// Print 日志写入缓冲管道
func (log *Logger) Print(message string) {
	logChannel <- map[string]interface{}{
		"msg": message,
	}
}

// Print 外部调用
func Print(message string) {
	logChannel <- map[string]interface{}{
		"msg": message,
	}
}

// New 自定义日志
func New() *Logger {
	l := &Logger{}
	return l
}

// Writer log writer interface
type Writer interface {
	Printf(string, ...interface{})
}

// Printf gorm需要的自定义日志接口
func (log *Logger) Printf(format string, v ...any) {
	format = strings.Replace(format, "\n", " ", -1)
	message := fmt.Sprintf(format, v...)
	log.Print(join("INFO", message))
}
