package master

import (
	"context"
	"sync"
	"time"

	"apollo/common"
	"apollo/utils"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

var log = utils.Logger
// master election connect etcd
type ElectionLock struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
	isLeader bool //是否选举称为leader
}
var (
	G_electionMaster *ElectionLock
)
func InitElectionMaster()(err error){
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)
	//初始化配置
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,                                     //集群地址
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond, //连接超时时间
	}

	//建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}
	//获取KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	G_electionMaster = &ElectionLock{
		client:  client,
		kv:      kv,
		lease:   lease,
	}


	return
}

/**
election master
 */
func Campaign(client *clientv3.Client, parentCtx context.Context, wg *sync.WaitGroup)(result <-chan struct{}){
	ip , err := utils.GetLocalIP()
	if err != nil{
		log.Warn(common.ERR_NO_LOCAL_IP_FOUND)
	}
	ctx,_ := context.WithCancel(parentCtx)

	//通知外层协程是否结束
	if wg != nil{
		wg.Add(1)
	}
	// 信号channel, 所有节点监听channel,使阻塞的几点可以成为leader，避免轮询是否是leader节点
	// 返回只读channel,所有节点可以阻塞
	notify := make(chan struct{},100)

	go func(){
		defer func(){
			if wg != nil{
				wg.Done()
			}
		}()
		for {
			select {
			case <- ctx.Done():
				return
			default:
			}

			session ,err := concurrency.NewSession(client,concurrency.WithTTL(5))
			// 生成session失败后两秒后重试
			if err != nil{
				time.Sleep(time.Second * 2)
				continue
			}

			//创建新的etcd选举election
			election := concurrency.NewElection(session,common.JOB_ELECTION_MASTER)

			//调用Campaign方法进行选举，leader节点会选举出来，非leader节点会阻塞在里面
			if err = election.Campaign(ctx, ip); err != nil{
				//选举失败重试
				//1.关闭旧的session
				if err = session.Close();err != nil{
					log.Warn("session close failed",ip)
				}
				time.Sleep(1 * time.Second)
				continue
			}

			breakFlag := false
			for !breakFlag{
				select {
				case notify <- struct{}{}: //向其余节点发信号
				case <-session.Done(): //如果因为网络原因导致etcd断开keepalive,break退出，重新选举
					breakFlag = true
					break
				case <- ctx.Done():
					ctxTemp , _ := context.WithTimeout(context.Background(), 1 * time.Second)
					//放弃leader角色
					err = election.Resign(ctxTemp)
					if err != nil{
						log.Warn("leader give up failed ",ip)
					}
					//关闭session
					err = session.Close()
					if err != nil{
						log.Warn("close old session failed",ip)
					}
					return
				}

			}
		}


	}()
	return notify
}
