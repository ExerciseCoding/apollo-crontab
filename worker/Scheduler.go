package worker

import "crontab/common"

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
//调度协程
func(scheduler *Scheduler) schedulerLoop(){
	//定时任务common.CronJob
	var(
		jobEvent *common.JobEvent
	)
	for{
		select {
		case jobEvent = <- scheduler.jobEventChan: //监听任务变化事件
		//对任务列表做操作
		}
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
