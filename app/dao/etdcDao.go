package dao

import (
	"github.com/coreos/etcd/clientv3"
	"zhyu/utils"
	"zhyu/utils/logger"
)

// Put 插入
func Put(key string, val string) {
	e := utils.EtcdObject
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
		logger.Warn("putResp.Header.Revision:%d", putResp.Header.Revision)
		//fmt.Println(putResp.Header.Revision) //输出Revision
		//fmt.Println(putResp.PrevKv.Value)    //输出修改前的值
	}
}

// Get 取出
func Get(key string) (val map[string]any) {
	e := utils.EtcdObject
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
	} else {
		//fmt.Println(getResp.Kvs)
	}
	//遍历键值对
	val = make(map[string]any)
	for _, ev := range getResp.Kvs {
		//fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		val[string(ev.Key)] = string(ev.Value)
	}
	return
}
