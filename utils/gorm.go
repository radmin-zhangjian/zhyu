package utils

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
	"zhyu/setting"
	logs "zhyu/utils/logger"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func InitDB() {
	// 写日志&控制台输出
	//filePath := logs.FilePath()
	//file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	//multiOutput := io.MultiWriter(os.Stdout, file)
	//multiLogger := log.New(multiOutput, "["+setting.Server.ServerName+"]", log.LstdFlags)
	// 只控制台输出
	//file, err := os.Create("gorm-log.txt")
	//fileLogger := log.New(file, "", log.LstdFlags)
	// 自定义log
	customLogger := logs.New()
	// 初始化数据库日志
	newLogger := logger.New(
		//log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		//multiLogger,
		customLogger,
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	fmt.Println(mySQLUri())
	conn, err1 := gorm.Open(mysql.Open(mySQLUri()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "zhyu_",
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err1 != nil {
		log.Printf("connect get failed.")
		return
	}
	sqlDB, err := conn.DB()
	if err != nil {
		log.Printf("database setup error %v", err)
	}
	sqlDB.SetMaxIdleConns(int(setting.Database.MaxIdleConn))                                //最大空闲连接数
	sqlDB.SetMaxOpenConns(int(setting.Database.MaxOpenConn))                                //最大连接数
	sqlDB.SetConnMaxLifetime(time.Duration(setting.Database.ConnMaxLifetime) * time.Second) //设置连接空闲超时
	db = conn
}

// 获取链接URI
func mySQLUri() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		setting.Database.UserName,
		setting.Database.Password,
		setting.Database.Host,
		setting.Database.Port,
		setting.Database.DbName)
}
