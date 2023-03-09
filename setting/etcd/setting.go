package etcd

import (
	"encoding/json"
	"flag"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/storage/storagepb"
	//"github.com/coreos/etcd/clientv3"
	//"github.com/coreos/etcd/storage/storagepb"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"zhyu/utils"
	"zhyu/utils/logger"
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
	utils.GetEtcd().Put("/setting/config/server", string(serverJson))
	databaseJson, _ := json.Marshal(Database)
	utils.GetEtcd().Put("/setting/config/db", string(databaseJson))
	redisJson, _ := json.Marshal(Redis)
	utils.GetEtcd().Put("/setting/config/redis", string(redisJson))
	elasticJson, _ := json.Marshal(Elastic)
	utils.GetEtcd().Put("/setting/config/elastic", string(elasticJson))
	whiteIpListJson, _ := json.Marshal(WhiteList)
	utils.GetEtcd().Put("/setting/config/whiteIpList", string(whiteIpListJson))
}

// 读取配置
func getConf() {
	// 监控所有配置文件
	s := &settingWatch{etcd: utils.GetEtcd(), key: "/setting/config/"}
	go s.listenWatchConfigs()

	// 单个配置监听
	//e := &settingWatch{etcd: utils.GetEtcd(), key: "/setting/config/server"}
	//go ListenWatch(e)

	respServer := utils.GetEtcd().Get("/setting/config/server")
	serverMap := server{}
	json.Unmarshal([]byte(respServer["/setting/config/server"].(string)), &serverMap)
	fmt.Println("serverJson:====================\n", serverMap)

	resDb := utils.GetEtcd().Get("/setting/config/db")
	databaseMap := database{}
	json.Unmarshal([]byte(resDb["/setting/config/db"].(string)), &databaseMap)
	fmt.Println("databaseJson:====================\n", databaseMap)

	resRedis := utils.GetEtcd().Get("/setting/config/redis")
	redisMap := redis{}
	json.Unmarshal([]byte(resRedis["/setting/config/redis"].(string)), &redisMap)
	fmt.Println("redisJson:====================\n", redisMap)

	resElastic := utils.GetEtcd().Get("/setting/config/elastic")
	elasticMap := elastic{}
	json.Unmarshal([]byte(resElastic["/setting/config/elastic"].(string)), &elasticMap)
	fmt.Println("elasticJson:====================\n", elasticMap)

	resWhiteIpList := utils.GetEtcd().Get("/setting/config/whiteIpList")
	ipMap := whiteList{}
	json.Unmarshal([]byte(resWhiteIpList["/setting/config/whiteIpList"].(string)), &ipMap)
	fmt.Println("whileIpJson:====================\n", ipMap)

	Server = &serverMap
	Database = &databaseMap
	Redis = &redisMap
	Elastic = &elasticMap
	WhiteList = &ipMap
}

type etcdWatch interface {
	listenWatchConfig()
}

func ListenWatch(s etcdWatch) {
	s.listenWatchConfig()
}

type settingWatch struct {
	etcd *utils.EtcdContext
	key  string
}

func (w *settingWatch) listenWatchConfig() {
	cli, err := clientv3.New(w.etcd.Config)
	if err != nil {
		logger.Error("connect to etcd failed, err:%v", err)
		return
	}
	defer cli.Close()

	// watch 监听
	rch := cli.Watch(w.etcd.Ctx, w.key) // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			if string(ev.Kv.Key) == w.key {
				serverMap := server{}
				json.Unmarshal([]byte(string(ev.Kv.Value)), &serverMap)
				Server = &serverMap
				logger.Info("WatchServer:%v", serverMap)
			}
		}
	}
}

// 监控所有配置
func (w *settingWatch) listenWatchConfigs() {
	cli, err := clientv3.New(w.etcd.Config)
	if err != nil {
		logger.Error("connect to etcd failed, err:%v", err)
		return
	}
	defer cli.Close()

	// watch 监听
	rch := cli.Watch(w.etcd.Ctx, w.key, clientv3.WithPrefix()) // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.Event_EventType(storagepb.PUT):
				if string(ev.Kv.Key) == "/setting/config/server" {
					serverMap := server{}
					json.Unmarshal([]byte(string(ev.Kv.Value)), &serverMap)
					Server = &serverMap
					logger.Info("WatchServerConfig:%v", serverMap)
				}
				if string(ev.Kv.Key) == "/setting/config/db" {
					dbMap := database{}
					json.Unmarshal([]byte(string(ev.Kv.Value)), &dbMap)
					Database = &dbMap
					logger.Info("WatchDatabaseConfig:%v", dbMap)
				}
				if string(ev.Kv.Key) == "/setting/config/redis" {
					redisMap := redis{}
					json.Unmarshal([]byte(string(ev.Kv.Value)), &redisMap)
					Redis = &redisMap
					logger.Info("WatchRedisConfig:%v", redisMap)
				}
				if string(ev.Kv.Key) == "/setting/config/elastic" {
					elasticMap := elastic{}
					json.Unmarshal([]byte(string(ev.Kv.Value)), &elasticMap)
					Elastic = &elasticMap
					logger.Info("WatchElasticConfig:%v", elasticMap)
				}
				if string(ev.Kv.Key) == "/setting/config/whiteIpList" {
					whiteIpListMap := whiteList{}
					json.Unmarshal([]byte(string(ev.Kv.Value)), &whiteIpListMap)
					WhiteList = &whiteIpListMap
					logger.Info("WatchWhiteIpListConfig:%v", whiteIpListMap)
				}
			case mvccpb.Event_EventType(storagepb.DELETE):
				logger.Info("删除:", "Revision: %v", ev.Kv.ModRevision)
			}
		}
	}
}
