package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"zhyu/app"
	"zhyu/app/common"
	"zhyu/app/dao"
	"zhyu/app/model/esModel"
	"zhyu/utils"
)

var tarIndex = "test_active_record_2023"

// EsCreateService 插入数据
func EsCreateService(c *app.Context) any {
	// 拼装数据
	item := esModel.TestActiveRecordModel{
		Controller: "login",
		Action:     "login",
		ShopCode:   "01133466",
		Ip:         "127.0.0.1",
		UserInfo: esModel.UserInfo{
			Id:       "1000",
			UserName: "gogo",
			Phone:    "9911199990",
			Address:  "北京CBD",
		},
		Datetime:  "2023-02-27 21:47:56",
		Timestamp: "1677505676"}

	// 转化成map
	docData := utils.StructToMap(item)
	c.Logs.Info("emJson docData", docData)

	// 写入ES
	var code int
	if ok := dao.EsCreateDao(tarIndex, docData); ok {
		code = common.SUCCESS
	} else {
		code = common.ERROR
	}

	return common.Result(code, common.GetMsg(code), nil)
}

// EsUpdateService 修改数据
func EsUpdateService(c *app.Context) any {
	// 拼装数据
	item := esModel.TestActiveRecordModel{
		Controller: "login",
		Action:     "sms",
		ShopCode:   "01133466",
		Ip:         "127.0.0.1",
		UserInfo: esModel.UserInfo{
			Id:       "1000",
			UserName: "gogo",
			Phone:    "9911199990",
			Address:  "北京CBD",
		},
		Datetime:  "2023-02-27 21:47:56",
		Timestamp: "1677505676"}

	// 转化成map
	docData := utils.StructToMap(item)
	c.Logs.Info("emJson docData", docData)

	// 写入ES
	var code int
	if ok := dao.EsUpdateDao(tarIndex, "ZnnsiIYBOKp382IudzuS", docData); ok {
		code = common.SUCCESS
	} else {
		code = common.ERROR
	}

	return common.Result(code, common.GetMsg(code), nil)
}

// EsGetService 查找单个数据
func EsGetService(c *app.Context) any {
	// 查询单个  result 是一个*json数据
	result := dao.EsGetDao(tarIndex, "ZnnsiIYBOKp382IudzuS")
	//c.Logs.Info("esGet result", result)
	var code int
	code = common.SUCCESS

	return common.Result(code, common.GetMsg(code), result)
}

// EsSearchService 查找多个数据
func EsSearchService(c *app.Context) any {
	// 查询
	resp := dao.EsSearchDao(tarIndex, "shopCode", "01133466", 10, 1)
	//c.Logs.Info("EsSearch result", result)

	// 解析 Hits 数据
	total := resp.Hits.TotalHits.Value // 文档总数
	var list []map[string]interface{}
	if total > 0 {
		for _, hit := range resp.Hits.Hits {
			var t map[string]interface{}
			err := json.Unmarshal(hit.Source, &t)
			if err != nil {
				fmt.Println("failed", err)
			}
			list = append(list, t)
		}
	}
	fmt.Printf("Hits.Hits: %v \n", list)

	// 组合数据
	var result = make(map[string]interface{})
	result["total"] = total
	result["list"] = list

	var code int
	code = common.SUCCESS

	return common.Result(code, common.GetMsg(code), result)
}

// EsQueryService 查找多个数据
func EsQueryService(c *app.Context) any {
	// query 条件
	queryStr := `{
		"size": %d,
		"query": {
			"bool": {
				"must": [
					{
						"term": {
							"shopCode": "%s"
						}
					},
					{
						"term": {
							"controller": "%s"
						}
					}
				]
			}
		},
		"aggs": {
			"groupByShopCode": {
			  "terms": {
				"field": "shopCode",
				"size": 10
			  },
			  "aggs": {
				"groupByDay": {
				  "date_histogram": {
					"field": "datetime",
					"calendar_interval": "day",
					"format": "yyyy-MM-dd",
					"min_doc_count": 1
				  },
				  "aggs": {
					"cardinalityDay": {
					  "cardinality": {
						"field": "datetime",
						"precision_threshold": 40000
					  }
					}
				  }
				}
			  }
			}
		}
	}`
	queryStr = fmt.Sprintf(queryStr, 10, "01133466", "login")
	// 查询
	resp := dao.EsQueryDao(tarIndex, queryStr)
	//c.Logs.Info("EsSearch resp info:", resp)

	// 解析 Hits 数据
	total := resp.Hits.TotalHits.Value // 文档总数
	var list []map[string]interface{}
	if total > 0 {
		for _, hit := range resp.Hits.Hits {
			var t map[string]interface{}
			err := json.Unmarshal(hit.Source, &t)
			if err != nil {
				fmt.Println("failed", err)
			}
			list = append(list, t)
		}
	}
	fmt.Printf("Hits.Hits: %v \n", list)

	// 自定义解析 aggregations 桶数据
	// 直接转json 这需要各个参数和进行转化相应类型才能使用
	var aggregationsMap map[string]interface{}
	emJson, _ := json.Marshal(resp.Aggregations)
	json.Unmarshal(emJson, &aggregationsMap)
	var aggsList []map[string]interface{}
	aggs1 := aggregationsMap["groupByShopCode"].(map[string]interface{})["buckets"].([]interface{})
	for _, val := range aggs1 {
		//fmt.Printf("AggregationsMap val: %v \n", val)
		aggs1Map := make(map[string]interface{})
		aggs1Map["shopCode"] = val.(map[string]interface{})["key"]
		aggs1Map["ret1Count"] = val.(map[string]interface{})["doc_count"]
		aggs2 := val.(map[string]interface{})["groupByDay"].(map[string]interface{})["buckets"].([]interface{})
		var aggsArr []map[string]interface{}
		for _, val2 := range aggs2 {
			//fmt.Printf("AggregationsMap aggs2: %v \n", val2)
			count := val2.(map[string]interface{})["cardinalityDay"].(map[string]interface{})["value"]
			fmt.Printf("AggregationsMap aggs count: %v \n", count)
			aggs2Map := make(map[string]interface{})
			aggs2Map["day"] = val2.(map[string]interface{})["key_as_string"]
			aggs2Map["ret2Count"] = val2.(map[string]interface{})["doc_count"]
			aggs2Map["ret3Count"] = count
			aggsArr = append(aggsArr, aggs2Map)
		}
		aggs1Map["list"] = aggsArr
		aggsList = append(aggsList, aggs1Map)
	}
	fmt.Printf("Aggregations aggs1: %v \n", aggsList)

	// 多级聚合查询 三方插件解析
	// 解析 aggregations 桶数据
	ret, ok := resp.Aggregations.Terms("groupByShopCode")
	if !ok {
		log.Println("agg results is nil")
	}
	var bucketList = make([]map[string]interface{}, len(ret.Buckets))
	for i, v := range ret.Buckets {
		bucketMap := make(map[string]interface{})
		bucketMap["shopCode"] = ret.Buckets[i].Key.(string)
		bucketMap["ret1Count"] = ret.Buckets[i].DocCount
		ret2, _ := v.Aggregations.Terms("groupByDay")
		for j, t := range ret2.Buckets {
			bucketMap["day"] = t.KeyAsString
			bucketMap["ret2Count"] = t.DocCount
			ret3, _ := ret2.Buckets[j].Aggregations.Terms("cardinalityDay")
			bucketMap["ret3Count"] = ret3.Aggregations["value"]
		}
		bucketList[i] = bucketMap
	}
	fmt.Printf("Aggregations bucketKeys: %v \n", bucketList)

	// 组合数据
	var result = make(map[string]interface{})
	result["total"] = total                  // hist文档总数
	result["list"] = list                    // hist文档数据
	result["aggregations"] = aggregationsMap // aggs原始map
	result["aggsList"] = aggsList            // aggs自定义解析的数据
	result["bucketList"] = bucketList        // aggs根据三方插件解析的数据

	// code
	var code int
	code = common.SUCCESS

	return common.Result(code, common.GetMsg(code), result)
}
