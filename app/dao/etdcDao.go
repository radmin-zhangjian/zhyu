package dao

import (
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"golang.org/x/net/context"
	"time"
	"zhyu/utils"
	"zhyu/utils/logger"
)

// Put 插入
func Put(key string, val string) {
	e := utils.GetEtcd()
	//建立连接
	cli, err := clientv3.New(e.Config)
	if err != nil {
		logger.Error("connect to etcd failed, err:%v", err)
		return
	}
	defer cli.Close()

	//写etcd中的键值对
	kv := clientv3.NewKV(cli)
	//putResp, err = kv.Put(context.Background(), "/setting/config", "ok", clientv3.WithPrevKV())
	putResp, err := kv.Put(e.Ctx, key, val)
	if err != nil {
		logger.Warn("connect to etcd failed, err:%v", err)
	} else {
		logger.Info("putResp.Header.Revision:%d", putResp.Header.Revision)
		//fmt.Println(putResp.Header.Revision) //输出Revision
		//fmt.Println(putResp.PrevKv.Value)    //输出修改前的值
	}
}

// Get 取出
func Get(key string) (val map[string]any) {
	e := utils.GetEtcd()
	//建立连接
	cli, err := clientv3.New(e.Config)
	if err != nil {
		logger.Error("connect to etcd failed, err:%v", err)
		return
	}
	defer cli.Close()

	//context超时控制
	//ctx, cancel := context.WithTimeout(e.Ctx, 1*time.Second)
	//resp, err := cli.Get(ctx, key)
	//cancel()
	//if err != nil {
	//	logger.Error("connect to etcd failed, err:%v", err)
	//}
	////遍历键值对
	//val = make(map[string]any)
	//for _, ev := range resp.Kvs {
	//	//fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	//	val[string(ev.Key)] = string(ev.Value)
	//}

	//读取etcd的键值
	kv := clientv3.NewKV(cli)
	getResp, err := kv.Get(e.Ctx, key)
	if err != nil {
		logger.Warn("connect to etcd failed, err:%v", err)
	}
	//遍历键值对
	val = make(map[string]any)
	for _, ev := range getResp.Kvs {
		//fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		val[string(ev.Key)] = string(ev.Value)
	}
	return
}

// EtcdLock 分布式锁
// m 锁对象  cancel 关闭连接的闭包
func EtcdLock(lockKey string) (m *concurrency.Mutex, cancel func()) {
	e := utils.GetEtcd()
	cli, err := clientv3.New(e.Config)
	if err != nil {
		logger.Error("connect to etcd failed, err:%v", err)
		return
	}

	if lockKey == "" || lockKey == "nil" {
		lockKey = "/lock"
	}
	//fmt.Println("lockKey:", lockKey)

	// 获取锁对象
	session, err := concurrency.NewSession(cli)
	if err != nil {
		cli.Close()
		logger.Error("concurrency to session, err:%v", err)
	}
	m = concurrency.NewMutex(session, lockKey)

	// 关闭连接
	cancel = func() {
		cli.Close()
		//fmt.Println("======== lockKey close ========")
	}

	return m, cancel
}

// EtcdTransaction 事务
func EtcdTransaction() {
	e := utils.GetEtcd()
	client, err := clientv3.New(e.Config)
	if err != nil {
		logger.Error("connect to etcd failed, err:%v", err)
		return
	}
	defer client.Close()

	//1. 上锁，创建租约
	lease := clientv3.NewLease(client)

	//申请一个5秒的租约
	leaseGrantResp, err := lease.Grant(e.Ctx, 20)
	if err != nil {
		fmt.Println(err)
		return
	}

	//拿到租约ID
	leaseId := clientv3.LeaseID(leaseGrantResp.ID)

	//取消续租
	ctx, cancel := context.WithCancel(e.Ctx)

	//确保函数推出后，自动续租会停止
	defer cancel()
	defer lease.Revoke(e.Ctx, leaseId)

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
				if keepResChan == nil {
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
	txn := kv.Txn(e.Ctx)

	//定义事务
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/1"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/1", "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/1")) //否则抢锁失败

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

}
