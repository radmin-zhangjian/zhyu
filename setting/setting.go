package setting

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
)

func init() {
	InitConf("configFile/app.yaml")
}

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
