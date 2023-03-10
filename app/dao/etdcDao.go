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

// NewEtcdConnect 获取etcd连接
func NewEtcdConnect() (txn *clientv3.Client, cancel func()) {
	e := utils.GetEtcd()
	cli, err := clientv3.New(e.Config)
	if err != nil {
		logger.Error("connect to etcd failed, err:%v", err)
		return
	}

	// 关闭连接
	cancel = func() {
		cli.Close()
		fmt.Println("======== etcd connect close ========")
	}

	return cli, cancel
}

// Put 插入
// opts 租约时间设置
func Put(key string, val string, opts ...int64) {
	e := utils.GetEtcd()
	//建立连接
	cli, err := clientv3.New(e.Config)
	if err != nil {
		logger.Error("connect to etcd failed, err:%v", err)
		return
	}
	defer cli.Close()

	// 输出修改前的值
	optsV3 := []clientv3.OpOption{
		clientv3.WithPrevKV(),
	}

	var leaseId clientv3.LeaseID
	for opt := range opts {
		// 创建租约
		leaseGrantResp, err := cli.Grant(context.TODO(), int64(opt))
		if err != nil {
			logger.Warn("crant create err:%v", err)
		}
		leaseId = leaseGrantResp.ID
		// 设置租约
		optsV3 = append(optsV3, clientv3.WithLease(leaseId))
	}

	//写etcd中的键值对
	kv := clientv3.NewKV(cli)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	putResp, err := kv.Put(ctx, key, val, optsV3...)
	cancel()
	if err != nil {
		logger.Warn("connect to etcd failed, err:%v", err)
	} else {
		logger.Info("putResp.Header.Revision:%d", putResp.Header.Revision) //输出Revision
		//logger.Info("putResp.PrevKv.Value:%v", string(putResp.PrevKv.Value)) //输出修改前的值
	}

	// op方式操作
	//ops := []clientv3.Op{
	//	clientv3.OpPut("aaa", "123"),
	//	clientv3.OpGet("aaa"),
	//}
	//for _, op := range ops {
	//	resp, err := cli.Do(context.TODO(), op)
	//	if  err != nil {
	//		log.Fatal(err)
	//	}
	//	if op.IsPut() {
	//		fmt.Println(resp.Put())
	//	}
	//	if op.IsGet()  {
	//		fmt.Println(resp.Get())
	//	}
	//}
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

	//读取etcd的键值
	kv := clientv3.NewKV(cli)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	getResp, err := kv.Get(ctx, key)
	cancel()
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
	//session, err := concurrency.NewSession(cli, concurrency.WithTTL(20))  // 自定义过期时间
	session, err := concurrency.NewSession(cli) // 默认过期时间60秒
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

// EtcdLease 租约&续租
func EtcdLease() {
	e := utils.GetEtcd()
	client, err := clientv3.New(e.Config)
	if err != nil {
		logger.Error("connect to etcd failed, err:%v", err)
		return
	}
	defer client.Close()

	//1. 上锁，创建租约
	lease := clientv3.NewLease(client)

	//申请一个20秒的租约
	leaseGrantResp, err := lease.Grant(e.Ctx, 20)
	if err != nil {
		fmt.Println(err)
		return
	}

	//拿到租约ID
	leaseId := leaseGrantResp.ID

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

}
