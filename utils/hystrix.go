package utils

import (
	"context"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"sync"
)

var hystrixObject *HystrixContext
var OnceH sync.Once

// HystrixContext 结构体
type HystrixContext struct {
	Ctx context.Context
}

// NewHystrix 构造对象
func NewHystrix() *HystrixContext {
	OnceH.Do(func() {
		hystrixObject = &HystrixContext{Ctx: context.Background()}
	})
	return hystrixObject
}

// SetTimeout 执行的超时时间
func (h *HystrixContext) SetTimeout(key string, val int) *HystrixContext {
	hystrix.ConfigureCommand(key, hystrix.CommandConfig{
		Timeout: val,
	})
	return hystrixObject
}

// SetMaxConcurrentRequests 最大并发量
func (h *HystrixContext) SetMaxConcurrentRequests(key string, val int) *HystrixContext {
	hystrix.ConfigureCommand(key, hystrix.CommandConfig{
		MaxConcurrentRequests: val,
	})
	return hystrixObject
}

// SetMaxRequestVolumeThreshold
// 一个统计窗口 10 秒内请求数量
// 达到这个请求数量后才去判断是否要开启熔断
func (h *HystrixContext) SetMaxRequestVolumeThreshold(key string, val int) *HystrixContext {
	hystrix.ConfigureCommand(key, hystrix.CommandConfig{
		RequestVolumeThreshold: val,
	})
	return hystrixObject
}

// SetSleepWindow
// 熔断器被打开后
// SleepWindow 的时间就是控制过多久后去尝试服务是否可用了
// 单位为毫秒
func (h *HystrixContext) SetSleepWindow(key string, val int) *HystrixContext {
	hystrix.ConfigureCommand(key, hystrix.CommandConfig{
		SleepWindow: val,
	})
	return hystrixObject
}

// SetErrorPercentThreshold
// 错误百分比
// 请求数量大于等于 RequestVolumeThreshold 并且错误率到达这个百分比后就会启动熔断
func (h *HystrixContext) SetErrorPercentThreshold(key string, val int) *HystrixContext {
	hystrix.ConfigureCommand(key, hystrix.CommandConfig{
		ErrorPercentThreshold: val,
	})
	return hystrixObject
}

func (h *HystrixContext) Hystrix(serviceName string, f func() (any, error)) (any, error) {
	var body any
	err := hystrix.Do(serviceName, func() error {
		var err error
		body, err = f()
		return err
		//requestUrl := url.URL{
		//	Scheme:   "http",
		//	Host:     "127.0.0.1" + ":" + "9090",
		//	Path:     "/api/v1/test",
		//	RawQuery: "id=100",
		//}
		//resp, err := http.Get(requestUrl.String())
		//if err != nil {
		//	return err
		//}
		//body, err = ioutil.ReadAll(resp.Body)
		//jsonErr := json.Unmarshal(body, result)
		//if jsonErr != nil {
		//	return jsonErr
		//}
		//return nil
	}, func(e error) error {
		// 断路器打开时的处理逻辑，本示例是直接返回错误提示
		return errors.New("Http errors！")
	})

	if err == nil {
		return body, nil
	} else {
		return body, err
	}
}

func Hystrix(serviceName string, f func() (any, error)) (any, error) {
	return NewHystrix().Hystrix(serviceName, f)
}
