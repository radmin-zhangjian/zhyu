package toml

import (
	"github.com/pelletier/go-toml/v2"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	InitConf("/setting/configFile/app.toml")
}

func InitConf(dataFile string) {
	//_, filename, _, _ := runtime.Caller(0)
	//filePath := path.Join(path.Dir(filename), dataFile)
	rootPath, _ := os.Getwd() //获取项目根路径
	filePath := rootPath + dataFile

	_, err := os.Stat(filePath)
	if err != nil {
		log.Printf("common file path %s not exist", filePath)
	}

	// 方式一 速度快 但 需要绝对地址
	tomlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("tomlFile.Get err   #%v ", err)
	}
	t := &conf{}
	err = toml.Unmarshal(tomlFile, &t)
	if err != nil {
		log.Printf("Unmarshal: %v", err)
	}
	log.Printf("load conf toml success\n %v", t)
	log.Printf("t.DB.A.Host\n %v", t.DB.A.Host)

	// 方式二 稍慢 但 可以读取上传的文件
	//var (
	//	fp       *os.File
	//	fcontent []byte
	//)
	//t := &conf{}
	//if fp, err = os.Open(filePath); err != nil {
	//	fmt.Println("open error ", err)
	//	return
	//}
	//if fcontent, err = ioutil.ReadAll(fp); err != nil {
	//	fmt.Println("ReadAll error ", err)
	//	return
	//}
	//if err = toml.Unmarshal(fcontent, &t); err != nil {
	//	fmt.Println("toml.Unmarshal error ", err)
	//	return
	//}
	//log.Printf("load conf toml success\n %v", t)

	Server = &t.Srv
	Database = &t.DB
	Redis = &t.RedisConfig
	Elastic = &t.ES
	IpWhite = &t.IpWhite
}
