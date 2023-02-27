package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"log"
	"reflect"
	"zhyu/utils"
	"zhyu/utils/logger"
)

var tarType = "_doc"

// EsCreateDao 插入文档
func EsCreateDao(index string, doc any) bool {
	client := utils.GetES()

	//1.使用结构体方式存入到es里面
	//2.使用字符串方式存入到es里面
	//3.使用MAP字典方式存入到es里面
	resp, err := client.Index().Index(index).BodyJson(doc).Do(context.Background())
	if err != nil {
		logger.Error("es EsCreate error: %v", err)
	}
	logger.Info("EsCreate: indexed %d to index %s, type %s, result %s", resp.Id, resp.Index, resp.Type, resp.Result)
	if resp.Result == "created" || resp.Result == "updated" {
		return true
	}
	return false
}

// EsUpdateDao 修改文档
// 注意用update更新个别字段时，请用map结构，struct会更新全部的。
func EsUpdateDao(index string, id string, doc any) bool {
	client := utils.GetES()

	//1.使用结构体方式存入到es里面
	//2.使用字符串方式存入到es里面
	//3.使用MAP字典方式存入到es里面
	resp, err := client.Update().Index(index).Id(id).
		Doc(doc).DocAsUpsert(true).Do(context.Background())
	if err != nil {
		logger.Error("es EsUpdate error: %v", err)
	}
	logger.Info("EsUpdate: indexed %d to index %s, type %s, result %s", resp.Id, resp.Index, resp.Type, resp.Result)
	if resp.Result == "created" || resp.Result == "updated" || resp.Result == "noop" {
		return true
	}
	return false
}

// EsDelDao 删除文档
func EsDelDao(index string, docId string) any {
	client := utils.GetES()

	resp, err := client.Delete().Index(index).Id(docId).Do(context.Background())
	if err != nil {
		logger.Error("es EsDel error: %v", err)
		return false
	}
	if resp.Result == "deleted" {
		logger.Error("es EsDel resp.Result: %v", resp.Result)
	}
	return true
}

// EsBulkInsertDao 批量插入文档  bulkNum批量写入的数量
func EsBulkInsertDao(index string, docs []interface{}, bulkNum int) bool {
	client := utils.GetES()

	length := len(docs)
	// case-1, no need to care add, just add and do
	if length < bulkNum {
		bulkReq := client.Bulk()
		for i := 0; i < length; i++ {
			doc := docs[i]
			req := elastic.NewBulkIndexRequest().Index(index).Type(tarType).Doc(doc)
			bulkReq = bulkReq.Add(req)
		}
		_, err := bulkReq.Do(context.Background())
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
	// case-2, length > bulkNum
	idx := 0
	for {
		bulkReq := client.Bulk()
		for i := 0; i < bulkNum; i++ {
			doc := docs[i]
			req := elastic.NewBulkIndexRequest().Index(index).Type(tarType).Doc(doc)
			bulkReq = bulkReq.Add(req)
			idx++
			if idx >= length {
				break
			}
		}
		_, err := bulkReq.Do(context.Background())
		if err != nil {
			log.Println(err)
			return false
		}
		if idx >= length {
			break
		}
	}
	return true
}

// EsGetDao 查找单个文档
func EsGetDao(index string, docId string) any {
	client := utils.GetES()

	resp, err := client.Get().Index(index).Id(docId).Do(context.Background())
	if err != nil {
		logger.Error("es EsGet error: %v", err)
	}

	if resp == nil {
		return false
	}

	//logger.Info("EsGet: indexed %d to index %s, type %s, found %s", resp.Id, resp.Index, resp.Type, resp.Found)

	return resp.Source // *json.RawMessage
}

// EsQueryDao 查找文档 || 桶查找
func EsQueryDao(index string, queryStr string) *elastic.SearchResult {
	client := utils.GetES()

	// 转换成map
	var sourceMap map[string]interface{}
	json.Unmarshal([]byte(queryStr), &sourceMap)

	// exec
	resp, err := client.Search().Index(index).Source(sourceMap).Do(context.Background())

	if err != nil {
		logger.Error("es EsQuery error: %v", err)
	}

	// resp 原始数据
	// resp.Hits  数据和总数
	// resp.Hits.Hits  只有数据没有总数
	return resp
}

// EsSearchDao 链式检索文档
func EsSearchDao(index string, name string, val interface{}, size int, page int) *elastic.SearchResult {
	client := utils.GetES()

	// 字符串模式 key:val 方式查询
	//query := elastic.NewQueryStringQuery("key:val")

	// 包含  about字段包含book
	//query := elastic.NewMatchPhraseQuery("about", "book")

	// 数据的bool查询 && 聚合查询 (通过构造NewxxxQuery)
	query := elastic.NewBoolQuery()
	if reflect.TypeOf(val).Kind() == reflect.Slice {
		query.Must(elastic.NewTermsQuery(name, val))
	} else {
		query.Must(elastic.NewTermsQuery(name, val))
	}

	// 范围查询
	//query.Must(elastic.NewRangeQuery(name).Gte(gte).Lte(lte)) // must_not/filter类似，故不再赘述

	// 分页 & 排序 === (Sort)
	//resp, err := client.Search().Index(index).Type(tarType).Query(query).
	//	Size(size).From((page-1)*size).Sort(field, ascend).Do(context.Background())

	// 复杂查询  所谓复杂，无非是多种条件加一起，通过接口提供的不断链式添加即可。
	// 记分数
	query.Must(
		elastic.NewTermQuery("controller", "login"),
		elastic.NewMatchPhraseQuery("action", "s"),
		//elastic.NewTermQuery("name-2", 100),
		//elastic.NewTermsQuery("name-3", []string{"val-1", "val-2"}),
		//elastic.NewRangeQuery("name-4").Gte(10).Lte(1000),
	)
	// 不计分数
	//query.Filter(
	//	elastic.NewTermQuery("name-5", "val"),
	//	elastic.NewTermQuery("name-6", 100),
	//	elastic.NewTermsQuery("name-7", []string{"val-1", "val-2"}),
	//)
	// must_not/filter like above

	// exec
	resp, err := client.Search().Index(index).Query(query).
		Size(size).From((page - 1) * size).Do(context.Background())

	if err != nil {
		logger.Error("es EsSearch error: %v", err)
	}

	//return resp.Hits.Hits
	return resp
}
