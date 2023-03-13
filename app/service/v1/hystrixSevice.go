package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"zhyu/app"
	"zhyu/app/common"
	"zhyu/utils"
)

func HystrixService(c *app.Context) any {
	// hystrix熔断用例
	fmt.Println("hystrix熔断用例===================")
	// 测试用例
	f := func() (interface{}, error) {
		requestUrl := url.URL{
			Scheme:   "http",
			Host:     "127.0.0.1" + ":" + "9090",
			Path:     "/api/v1/test",
			RawQuery: "id=100",
		}
		resp, err := http.Get(requestUrl.String())
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		return body, err
	}
	result, err := utils.Hystrix("Comments", f)
	// 返回数据
	var code int = common.SUCCESS
	if err != nil {
		code = common.ERROR
		result = []byte{}
	}
	return common.Result(code, common.GetMsg(code), string(result.([]byte)))
}
