package v1

import (
	"encoding/json"
	"fmt"
	"zhyu/app"
	"zhyu/app/common"
	"zhyu/app/dao"
	"zhyu/setting"
)

// EtcdPutService 插入数据
func EtcdPutService(c *app.Context) {
	dao.Put("/setting/config", "{\"aaa\":111,\"bbb\":222}")

	// 降级服务的简单应用
	fmt.Println("CallSecondaryService is ", common.CallSecondaryService)
	if common.CallSecondaryService != 0 {
		fmt.Println("我是非必要数据！！！")
	}
}

type database struct {
	Type            string `yaml:"type"`
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	UserName        string `yaml:"username"`
	Password        string `yaml:"password"`
	DbName          string `yaml:"dbname"`
	MaxIdleConn     int64  `yaml:"max_idle_conn"`
	MaxOpenConn     int64  `yaml:"max_open_conn"`
	ConnMaxLifetime int64  `yaml:"conn_max_lifetime"`
}

// EtcdGetService 取出数据
func EtcdGetService(c *app.Context) any {
	dbConfJson, _ := json.Marshal(setting.Database)
	fmt.Println("dbConfJson:\n", string(dbConfJson))
	resDb := dao.Get("/setting/config/db")
	dbMap := database{}
	json.Unmarshal([]byte(resDb["/setting/config/db"].(string)), &dbMap)
	fmt.Println("dbConfJson-db:\n", dbMap)
	return dbMap

	res := dao.Get("/setting/config")
	return res["/setting/config"].(string)
}
