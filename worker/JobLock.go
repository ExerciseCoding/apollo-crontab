package worker

import (
	"context"

	"github.com/coreos/etcd/clientv3"
)

type JobLock struct {
	kv clientv3.KV
	lease clientv3.Lease
	jobName string //任务名
	cancelFunc context.CancelFunc //用于终止自动续租

}

//初始化锁
func InitJobLock(jobName string,kv clientv3.KV, lease clientv3.Lease)(jobLock *JobLock){
	jobLock = &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
	return
}

//尝试上锁
func (jobLock *JobLock) TryLock(err error){
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx context.Context
		cancelFunc context.CancelFunc
		leaseId clientv3.LeaseID
		keepRespChan <- chan *clientv3.LeaseKeepAliveResponse
	)
	//1.创建租约(5秒)
	if leaseGrantResp,err = jobLock.lease.Grant(context.TODO(),5); err != nil{
		return
	}

	//context用于取消自动续租
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())
	//租约ID
	leaseId = leaseGrantResp.ID
	//2.自动续租
	if keepRespChan,err = jobLock.lease.KeepAlive(cancelCtx,leaseId); err != nil{
		goto FAIL
	}

	//处理续租应答的协程
	go func(){
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <- keepRespChan: //自动续租的应答
				if keepResp == nil{
					goto END
				}
			}
		}

		END:
	}()
	//3.创建事务txn

	//4.事务抢锁

	//5.成功返回，失败释放租约

FAIL:
	cancelFunc() //取消自动续租
	jobLock.lease.Revoke(context.TODO(),leaseId) //释放租约
	return
}