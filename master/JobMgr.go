package master

import (
	"context"
	"crontab/common"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

//定义任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

var (
	//单例
	G_jobMgr *JobMgr

)

//初始化管理器
func InitJobMgr()(err error){
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
	)
	//初始化配置
	config = clientv3.Config{
		Endpoints: G_config.EtcdEndpoints, //集群地址
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond ,//连接超时时间
	}

	//建立连接
	if client,err  = clientv3.New(config); err != nil{
		return
	}

	//获取KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	//赋值单例
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	fmt.Println(G_jobMgr)
	return
}

//保存任务接口
func (jobMgr *JobMgr) SaveJob(job *common.CronJob)(oldJob *common.CronJob,err error){
	//把任务保存到/cron/jobs/任务名->json
	var(
		jobValue []byte
		putResp *clientv3.PutResponse
		oldJobObj common.CronJob
	)
	//etcd的key
	jobKey := "/cron/jobs/" + job.Name
	//任务信息json
	if jobValue,err = json.Marshal(job); err != nil{
		return
	}

	//保存到etcd
	if putResp,err = jobMgr.kv.Put(context.TODO(),jobKey,string(jobValue),clientv3.WithPrevKV()); err != nil{
		return
	}

	//如果是更新，那么返回旧值
	if putResp.PrevKv != nil{
		//对旧值做一个反序列化
		if err = json.Unmarshal(putResp.PrevKv.Value,&oldJobObj); err != nil{
			err = nil
			return
		}
		oldJob = &oldJobObj
	}

	return
}