package v1

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"strconv"
	"time"
	"zhyu/app"
	"zhyu/app/common"
	"zhyu/app/dao"
)

// EtcdPutService 插入数据
func EtcdPutService(c *app.Context) {
	dao.Put("/setting/config", "{\"aaa\":111,\"bbb\":222}", 120)

	// 降级服务的简单应用
	fmt.Println("CallSecondaryService is ", common.CallSecondaryService)
	if common.CallSecondaryService != 0 {
		fmt.Println("我是非必要数据！！！")
	}
}

// EtcdGetService 取出数据
func EtcdGetService(c *app.Context) any {
	res := dao.Get("/setting/config")
	data := res["/setting/config"].(string)
	return common.Result(common.SUCCESS, common.GetMsg(common.SUCCESS), data)
}

// EtcdLockService 锁
func EtcdLockService(c *app.Context) {

	lockKey := "/lock"

	// 模拟锁请求
	go func() {
		// 获取锁
		m, cancel := dao.EtcdLock(lockKey)
		defer cancel()
		// 加锁
		if err := m.Lock(context.Background()); err != nil {
			log.Fatal("go1 get mutex failed " + err.Error())
		}
		fmt.Printf("go1 get mutex sucess\n")
		// 执行任务
		time.Sleep(time.Duration(2) * time.Second)
		// 解锁
		if err := m.Unlock(context.Background()); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("go1 release lock\n")
	}()

	m2Locked := make(chan struct{})
	go func() {
		defer close(m2Locked)
		// 获取锁
		m, cancel := dao.EtcdLock(lockKey)
		defer cancel()
		time.Sleep(time.Duration(1) * time.Second)
		// 加锁
		if err := m.Lock(context.Background()); err != nil {
			log.Fatal("go2 get mutex failed " + err.Error())
		}
		fmt.Printf("go2 get mutex sucess\n")
		// 执行任务
		time.Sleep(time.Duration(2) * time.Second)
		// 解锁
		if err := m.Unlock(context.Background()); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("go2 release lock\n")
	}()

	<-m2Locked
	fmt.Println("acquired lock for s2")
}

// EtcdTransactionService 事务测试
// 改变两个user的值
func EtcdTransactionService(c *app.Context) (result any) {
	// 返回参数
	result = common.Result(common.ERROR, common.GetMsg(common.ERROR), nil)

	// 连接etcd
	cli, cancel := dao.NewEtcdConnect()
	defer cancel() // 关闭连接

	type user struct {
		key         string
		account     int
		modRevision int64
	}
	userA := user{key: "userA"}
	userB := user{key: "userB"}

	kv := clientv3.NewKV(cli)
	// 初始化数据
	putUserA, _ := kv.Put(context.Background(), userA.key, strconv.Itoa(1000))
	fmt.Println("put userA:", putUserA)
	putUserB, _ := kv.Put(context.Background(), userB.key, strconv.Itoa(1000))
	fmt.Println("put userA:", putUserB)

	// 取数据
	getUserA, err := kv.Get(context.Background(), userA.key)
	if len(getUserA.Kvs) == 1 {
		userA.modRevision = getUserA.Kvs[0].ModRevision
		userA.account, _ = strconv.Atoi(string(getUserA.Kvs[0].Value))
	}
	fmt.Println("getUserA:", userA)
	getUserB, err := kv.Get(context.Background(), userB.key)
	if len(getUserA.Kvs) == 1 {
		userB.modRevision = getUserB.Kvs[0].ModRevision
		userB.account, _ = strconv.Atoi(string(getUserB.Kvs[0].Value))
	}
	fmt.Println("getUserB:", userB)

	userA.account -= 100
	userB.account += 100

	//创建事务
	txn := kv.Txn(context.Background())
	//定义事务  if比较结果为true，then设置它，else失败
	txn.If(
		clientv3.Compare(clientv3.ModRevision(userA.key), "=", userA.modRevision),
		clientv3.Compare(clientv3.ModRevision(userB.key), "=", userB.modRevision),
	).Then(
		clientv3.OpPut(userA.key, strconv.Itoa(userA.account)),
		clientv3.OpPut(userB.key, strconv.Itoa(userB.account)),
	).Else(
	//do something
	)

	// 提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		c.Logs.Info("commit 事务失败：", err)
		return
	}

	// 事务失败
	if !txnResp.Succeeded {
		c.Logs.Info("succeeded 事务失败：", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	// 成功返回
	fmt.Println("事务成功")
	data := ""
	result = common.Result(common.SUCCESS, common.GetMsg(common.SUCCESS), data)
	return
}

// EtcdLockTransactionService 事务锁测试
func EtcdLockTransactionService(c *app.Context) (result any) {
	// 返回参数
	result = common.Result(common.ERROR, common.GetMsg(common.ERROR), nil)

	// 连接etcd
	client, cancel := dao.NewEtcdConnect()
	defer cancel() // 关闭连接

	//1. 上锁，创建租约
	lease := clientv3.NewLease(client)

	//申请一个5秒的租约
	leaseGrantResp, err := lease.Grant(context.TODO(), 5)
	if err != nil {
		fmt.Println(err)
		return
	}

	//拿到租约ID
	leaseId := leaseGrantResp.ID

	//取消续租
	ctx, cancelFunc := context.WithCancel(context.TODO())

	//确保函数推出后，自动续租会停止
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

	//续租
	keepResChan, err := lease.KeepAlive(ctx, leaseId)
	if err != nil {
		fmt.Println(err)
		return
	}

	//处理续租应答的协程
	go func() {
		for {
			select {
			case keepRes := <-keepResChan:
				if keepResChan == nil || keepRes == nil {
					fmt.Println("租约已经失效")
					goto END
				} else {
					//每秒会续租一次，所以会收到一次应答
					fmt.Println("收到自动续租应答", keepRes.ID)
				}
			}
		}
	END:
	}()

	//if不存在key，then设置它，else抢锁失败
	kv := clientv3.NewKV(client)

	//创建事务
	txn := kv.Txn(context.TODO())

	//定义事务
	txn.If(
		clientv3.Compare(clientv3.CreateRevision("/cront/lock/1"), "=", 0),
	).Then(
		clientv3.OpPut("/cron/lock/1", "", clientv3.WithLease(leaseId)),
	).Else(
		clientv3.OpGet("/cron/lock/1"),
	) //否则抢锁失败

	//提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		fmt.Println(err)
		return
	}

	//判断是否抢到锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
	}

	//2.处理事务
	fmt.Println("处理事务")
	time.Sleep(5 * time.Second)

	//defer释放锁

	result = common.Result(common.SUCCESS, common.GetMsg(common.SUCCESS), nil)
	return
}
