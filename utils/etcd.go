package utils

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"strconv"
	"sync"
	"time"
	"zhyu/app/common"
	"zhyu/utils/logger"
)

// 启动 etcd
// etcd --listen-client-urls 'http://0.0.0.0:2379' --advertise-client-urls 'http://0.0.0.0:2379'
// 通过指定参数进行启动 etcd
// nohup ./etcd --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379 --listen-peer-urls http://0.0.0.0:2380 &
// 查看版本 etcd --version

var EtcdObject *EtcdClint
var Once sync.Once

// EtcdClint 结构体
type EtcdClint struct {
	Config clientv3.Config
	Ctx    context.Context
}

// 初始化
func init() {
	//NewEtcd()
}

// etcd配置信息
func confEtcd() (cf clientv3.Config) {
	cf = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	return
}

// NewEtcd 构造对象
func NewEtcd() *EtcdClint {
	Once.Do(func() {
		EtcdObject = &EtcdClint{Config: confEtcd(), Ctx: context.Background()}
	})
	return EtcdObject
}

// ListenReduceRank 监听降级服务
func (e *EtcdClint) ListenReduceRank() {
	cli, err := clientv3.New(e.Config)
	if err != nil {
		logger.Error("connect to etcd failed, err:%v", err)
		return
	}
	defer cli.Close()

	// watch 监听
	rch := cli.Watch(e.Ctx, common.SecondaryKey) // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			//fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			if string(ev.Kv.Key) == common.SecondaryKey {
				common.CallSecondaryService, _ = strconv.Atoi(string(ev.Kv.Value))
				//fmt.Printf("callCommentService:%d\n", common.CallSecondaryService)
				logger.Info("callCommentService:%d", common.CallSecondaryService)
			}
		}
	}
}

// Put 插入
func (e *EtcdClint) Put(key string, val string) {
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
		logger.Info("putResp.Header.Revision:%d", putResp.Header.Revision) //输出Revision
		//logger.Info(putResp.PrevKv.Value)    //输出修改前的值
	}
}

// Get 取出
func (e *EtcdClint) Get(key string) (val map[string]any) {
	//建立连接
	cli, err := clientv3.New(e.Config)
	if err != nil {
		logger.Error("connect to etcd failed, err:%v", err)
		return
	}
	defer cli.Close()

	//context超时控制
	//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//resp, err := cli.Get(ctx, "name")
	//cancel()
	//if err != nil {
	//	fmt.Printf("get from etcd failed,err %v\n", err)
	//}
	////遍历键值对
	//for _, kv := range resp.Kvs {
	//	fmt.Printf("%s:%s \n", kv.Key, kv.Value)
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
		val[string(ev.Key)] = string(ev.Value)
	}
	return
}
