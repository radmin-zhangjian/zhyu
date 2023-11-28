package setting

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"testing"
)

var (
	// 测试环境默认 configFile/app.yaml 为配置文件
	// 更改配置文件 go main.go -config.file=xxx.yaml
	configFile = flag.String("config.file", "setting/configFile/app.yaml", "config file")
)

func init() {
	// 单元测试时需要 不添加会报错
	testing.Init()
	// 解析参数
	if !flag.Parsed() {
		flag.Parse()
	}
	InitConf(*configFile)
}

// InitConf
// https://github.com/go-bindata/go-bindata
// go-bindata -o setting/bindata.go -pkg setting setting/configFile/...
// CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags '-w -s' -a -o app_zhyu .
func InitConf(dataFile string) {
	// 打包配置文件
	var err error
	yamlFile, err := Asset(dataFile)
	if err != nil {
		log.Fatalf("Fail to Asset parse 'setting/configFile/app.yaml': %v", err)
	}
	log.Printf("asset: %T %v\n", yamlFile, string(yamlFile))

	// yml 正常解析
	//_, filename, _, _ := runtime.Caller(0)
	//filePath := path.Join(path.Dir(filename), dataFile)
	//_, err := os.Stat(filePath)
	//if err != nil {
	//	log.Printf("common file path %s not exist", filePath)
	//}
	//yamlFile, err := ioutil.ReadFile(filePath)
	//if err != nil {
	//	log.Printf("yamlFile.Get err   #%v ", err)
	//}

	// 解析成 struct
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
}

func (c conf) String() string {
	return fmt.Sprintf("%v\n%v\n%v\n%v", c.Srv, c.DB, c.RedisConfig, c.ES)
}

func (s server) String() string {
	return fmt.Sprintf("server : \n"+
		"\tserverName : %v \n"+
		"\tport : %v \n"+
		"\trunMode : %v \n"+
		"\tlogLevel : %v \n"+
		"\tlogPath : %v \n"+
		"\treadTimeout : %v \n"+
		"\twriteTimeout : %v \n"+
		"\tshutdownTime : %v \n"+
		"\tworkerID : %v \n"+
		"\tjwtSecret : %v", s.ServerName, s.Port, s.RunMode, s.LogLevel, s.LogPath, s.ReadTimeout, s.WriteTimeout,
		s.ShutdownTime, s.WorkerID, s.JwtSecret)
}

func (m database) String() string {
	return fmt.Sprintf("database : \n"+
		"\ttype : %v \n"+
		"\thost : %v \n"+
		"\tport : %v \n"+
		"\tusername : %v \n"+
		"\tpassword : %v \n"+
		"\tdbname : %v \n"+
		"\tmax_idle_conn : %v \n"+
		"\tmax_open_conn : %v \n"+
		"\tconn_max_lifetime : %v",
		m.Type, m.Host, m.Port, m.UserName, m.Password, m.DbName, m.MaxOpenConn, m.MaxIdleConn, m.ConnMaxLifetime)
}
func (r redis) String() string {
	return fmt.Sprintf("redis : \n"+
		"\thost : %v \n"+
		"\tport : %v \n"+
		"\tpassword : %v \n"+
		"\tdb : %v \n"+
		"\tpoolSize : %v",
		r.Host, r.Port, r.Password, r.DB, r.PoolSize)
}
func (r elastic) String() string {
	return fmt.Sprintf("elastic : \n"+
		"\thost : %v \n"+
		"\tuser : %v \n"+
		"\tpassword : %v",
		r.Host, r.User, r.Password)
}
