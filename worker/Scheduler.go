package worker

import (
	"crontab/common"
	"fmt"
	"time"
)

//任务调度
type Scheduler struct {
	jobEventChan chan *common.JobEvent //etcd任务事件队列
	jobPlanTable map[string]*common.JobSchedulerPlan //任务调度计划表
}

var (
	G_scheduler *Scheduler
)

func (scheduler *Scheduler)handlerJobEvent(jobEvent *common.JobEvent){
	var(
		jobSchedulerPlan *common.JobSchedulerPlan
		jobExisted bool
		err error
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE: //保存任务时间
		if jobSchedulerPlan,err = common.BuildJobSchedulerPlan(jobEvent.Job); err != nil{
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulerPlan
	case common.JOB_EVENT_DELETE: //删除任务事件
		if jobSchedulerPlan,jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted{
			delete(scheduler.jobPlanTable,jobEvent.Job.Name)
		}
	}
}

func (scheduler *Scheduler) TryScheduler()(schedulerAfter time.Duration){
	var(
		jobPlan *common.JobSchedulerPlan
		now time.Time
		nearTime *time.Time
	)
	//如果任务为空
	if len(scheduler.jobPlanTable) == 0{
		schedulerAfter = 1 * time.Second
		return
	}
	now = time.Now()
	//遍历所有任务
	for _, jobPlan = range scheduler.jobPlanTable{
		if jobPlan.NextTime.Before(time.Now()) || jobPlan.NextTime.Equal(now){
			//TODO: 尝试执行任务
			fmt.Println("执行任务",jobPlan.Job.Name)
			jobPlan.NextTime = jobPlan.Expr.Next(now) //更新下次执行时间
		}

		//统计最近一个要过期的任务时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime){
			nearTime = &jobPlan.NextTime
		}

	}
	//下次调度间隔（最近要执行的任务调度时间-当前时间)
	schedulerAfter = (*nearTime).Sub(now)
	return
}
//调度协程
func(scheduler *Scheduler) schedulerLoop(){
	//定时任务common.CronJob
	var(
		jobEvent *common.JobEvent
		schedulerAfer time.Duration
		schedulerTimer *time.Timer
	)

	//初始化一次
	schedulerAfer  = scheduler.TryScheduler()

	//调度的延迟定时器
	schedulerTimer = time.NewTimer(schedulerAfer)
	for{
		select {
		case jobEvent = <- scheduler.jobEventChan: //监听任务变化事件
		//对任务列表做操作
			scheduler.handlerJobEvent(jobEvent)
		case <- schedulerTimer.C: //最近的任务到期
		}

		//调度一次任务
		schedulerAfer = scheduler.TryScheduler()
		//重置调度间隔
		schedulerTimer.Reset(schedulerAfer)
	}

}

//推送任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent){
	scheduler.jobEventChan <- jobEvent
}

//初始化调度器
func InitScheduler() (err error){
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.JobSchedulerPlan),
	}

	//启动调度协程
	go G_scheduler.schedulerLoop()
	return
}
