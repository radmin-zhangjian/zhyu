package etcd

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"zhyu/utils"
)

var (
	// 测试环境默认 configFile/app.yaml 为配置文件
	// 更改配置文件 go main.go -config.etcd.file=xxx.yaml
	configFile = flag.String("config.etcd.file", "../configFile/app.yaml", "config file")
)

func init() {
	// 解析参数
	if !flag.Parsed() {
		flag.Parse()
	}
	utils.NewEtcd()
	// 写配置 可以单独进行配置 这里只是测试用
	InitConf(*configFile)
	// 读取 etcd 配置
	getConf()
}

// InitConf 初始化配置
func InitConf(dataFile string) {
	_, filename, _, _ := runtime.Caller(0)
	filePath := path.Join(path.Dir(filename), dataFile)
	_, err := os.Stat(filePath)
	if err != nil {
		log.Printf("common file path %s not exist", filePath)
	}
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	c := &conf{}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Printf("Unmarshal: %v", err)
	}
	log.Printf("load conf success\n %v", c)

	Server = &c.Srv
	Database = &c.DB
	Redis = &c.RedisConfig
	Elastic = &c.ES
	WhiteList = &c.WhiteList

	// 配置写入配置中心
	serverJson, _ := json.Marshal(Server)
	utils.EtcdObject.Put("/setting/config/server", string(serverJson))
	databaseJson, _ := json.Marshal(Database)
	utils.EtcdObject.Put("/setting/config/db", string(databaseJson))
	redisJson, _ := json.Marshal(Redis)
	utils.EtcdObject.Put("/setting/config/redis", string(redisJson))
	elasticJson, _ := json.Marshal(Elastic)
	utils.EtcdObject.Put("/setting/config/elastic", string(elasticJson))
	whiteIpListJson, _ := json.Marshal(WhiteList)
	utils.EtcdObject.Put("/setting/config/whiteIpList", string(whiteIpListJson))
}

// 读取配置
func getConf() {
	respServer := utils.EtcdObject.Get("/setting/config/server")
	serverMap := server{}
	json.Unmarshal([]byte(respServer["/setting/config/server"].(string)), &serverMap)
	fmt.Println("serverJson:====================\n", serverMap)

	resDb := utils.EtcdObject.Get("/setting/config/db")
	databaseMap := database{}
	json.Unmarshal([]byte(resDb["/setting/config/db"].(string)), &databaseMap)
	fmt.Println("databaseJson:====================\n", databaseMap)

	resRedis := utils.EtcdObject.Get("/setting/config/redis")
	redisMap := redis{}
	json.Unmarshal([]byte(resRedis["/setting/config/redis"].(string)), &redisMap)
	fmt.Println("redisJson:====================\n", redisMap)

	resElastic := utils.EtcdObject.Get("/setting/config/elastic")
	elasticMap := elastic{}
	json.Unmarshal([]byte(resElastic["/setting/config/elastic"].(string)), &elasticMap)
	fmt.Println("elasticJson:====================\n", elasticMap)

	resWhiteIpList := utils.EtcdObject.Get("/setting/config/whiteIpList")
	ipMap := whiteList{}
	json.Unmarshal([]byte(resWhiteIpList["/setting/config/whiteIpList"].(string)), &ipMap)
	fmt.Println("whileIpJson:====================\n", ipMap)

	Server = &serverMap
	Database = &databaseMap
	Redis = &redisMap
	Elastic = &elasticMap
	WhiteList = &ipMap
}
