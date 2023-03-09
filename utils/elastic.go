package utils

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"zhyu/setting"
)

//var ESURL = "http://elastic:password@127.0.0.1:9200/"

var esClient *elastic.Client

type ElasticORM struct {
	client *elastic.Client
	index  string
	ctx    context.Context
}

func init() {
	//InitES()
}

func NewElasticORM() *ElasticORM {
	return &ElasticORM{
		client: GetES(),
		ctx:    context.Background(),
	}
}

// GetES 获取ES对象
func GetES() *elastic.Client {
	return esClient
}

// InitES 初始化es驱动
func InitES() {

	sniffOpt := elastic.SetSniff(false)                                             // 非集群下，关闭嗅探机制
	urlOpt := elastic.SetURL(setting.Elastic.Host)                                  // URL自行设置，比如 http://<user>:<passwd>@ip:9200
	authOpt := elastic.SetBasicAuth(setting.Elastic.User, setting.Elastic.Password) // 也可以把auth参数放到url中
	checkOpt := elastic.SetHealthcheck(false)
	errorLog := elastic.SetErrorLog(log.New(os.Stdout, "app", log.LstdFlags))
	options := []elastic.ClientOptionFunc{urlOpt, authOpt, sniffOpt, checkOpt, errorLog}

	client, err := elastic.NewClient(options...)
	if err != nil {
		log.Fatal(err)
	}
	// check es conn
	info, statusCode, err := client.Ping(setting.Elastic.Host).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if statusCode == 200 {
		log.Println("connect es success.")
	}
	fmt.Printf("Es return with code %d and version %s \n", statusCode, info.Version.Number)
	esClient = client

	//esversionCode, err := EsClient.ElasticsearchVersion(ESURL)
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Printf("es version %s\n", esversionCode)
	//return client
}
