package v1

import (
	"context"
	"fmt"
	"log"
	"time"
	"zhyu/app"
	"zhyu/app/common"
	"zhyu/app/dao"
)

// EtcdPutService 插入数据
func EtcdPutService(c *app.Context) {
	dao.Put("/setting/config", "{\"aaa\":111,\"bbb\":222}")

	// 降级服务的简单应用
	fmt.Println("CallSecondaryService is ", common.CallSecondaryService)
	if common.CallSecondaryService != 0 {
		fmt.Println("我是非必要数据！！！")
	}
}

// EtcdGetService 取出数据
func EtcdGetService(c *app.Context) any {
	res := dao.Get("/setting/config")
	return res["/setting/config"].(string)
}

// EtcdLockService 锁
func EtcdLockService(c *app.Context) {

	lockKey := "/lock"

	// 模拟锁请求
	go func() {
		// 获取锁
		m, cancel := dao.EtcdLock(lockKey)
		defer cancel()
		if err := m.Lock(context.Background()); err != nil {
			log.Fatal("go1 get mutex failed " + err.Error())
		}
		fmt.Printf("go1 get mutex sucess\n")
		fmt.Println(m)
		time.Sleep(time.Duration(10) * time.Second)
		if err := m.Unlock(context.Background()); err != nil {
			log.Fatal(err)
		}
		fmt.Println("released lock for s1")
		fmt.Printf("go1 release lock\n")
	}()

	m2Locked := make(chan struct{})
	go func() {
		defer close(m2Locked)
		// 获取锁
		m, cancel := dao.EtcdLock(lockKey)
		defer cancel()
		time.Sleep(time.Duration(2) * time.Second)
		if err := m.Lock(context.Background()); err != nil {
			log.Fatal("go2 get mutex failed " + err.Error())
		}
		fmt.Printf("go2 get mutex sucess\n")
		fmt.Println(m)
		time.Sleep(time.Duration(2) * time.Second)
		if err := m.Unlock(context.Background()); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("go2 release lock\n")
	}()

	<-m2Locked
	fmt.Println("acquired lock for s2")
}
