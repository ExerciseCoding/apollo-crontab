package worker

import (
	"fmt"
	"time"

	"apollo/common"
)

//任务调度
type Scheduler struct {
	jobEventChan chan *common.JobEvent //etcd任务事件队列
	jobPlanTable map[string]*common.JobSchedulerPlan //任务调度计划表
	jobExecutingTable map[string]*common.JobExecuteInfo //任务执行表
	jobResultChan chan *common.JobExecuteResult //任务结果队列
}

var (
	G_scheduler *Scheduler
)

func (scheduler *Scheduler)handlerJobEvent(jobEvent *common.JobEvent){
	var(
		jobSchedulerPlan *common.JobSchedulerPlan
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting bool
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
	case common.JOB_EVENT_KILL: //强杀任务事件
		//取消掉Command执行,判断任务是否在执行中
		if jobExecuteInfo,jobExecuting = scheduler.jobExecutingTable[jobEvent.Job.Name]; jobExecuting{
			jobExecuteInfo.CancelFunc() //触发command杀死shell子进程，任务得到退出
		}

	}
}
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulerPlan){
	//调度和执行是2件事情
	var(
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting bool
	)

	//执行的任务可能运行很久，1分钟会调度60次，但是只能执行1次,防止并发

	//如果任务正在执行，跳过本次调度
	if jobExecuteInfo,jobExecuting = scheduler.jobExecutingTable[jobPlan.Job.Name]; jobExecuting{
		//fmt.Println("尚未推出，跳过执行:",jobPlan.Job.Name)
		return
	}

	//构建执行状态
	jobExecuteInfo = common.BuildJobExcuteInfo(jobPlan)

	//保存执行状态
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	//fmt.Println("执行任务:",jobPlan.Job.Name,jobExecuteInfo.PlanTime,jobExecuteInfo.RealTime)
	//执行任务
	G_executor.ExecuteJob(jobExecuteInfo)


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
			scheduler.TryStartJob(jobPlan)
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
		jobResult *common.JobExecuteResult
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
		case jobResult = <- scheduler.jobResultChan: //监听任务执行结果
		scheduler.handlerJobResult(jobResult)
		}

		//调度一次任务
		schedulerAfer = scheduler.TryScheduler()
		//重置调度间隔
		schedulerTimer.Reset(schedulerAfer)
	}

}

//处理任务结果
func (scheduler *Scheduler) handlerJobResult(result * common.JobExecuteResult) {
	var (
		jobLog *common.JobLog
	)
	//删除执行状态
	delete(scheduler.jobExecutingTable, result.ExecuteInfo.Job.Name)

	//生成执行日志
	if result.Err != common.ERR_LOCK_ALREADY_REQUIRED {
		jobLog = &common.JobLog{
			JobName:      result.ExecuteInfo.Job.Name,
			Command:      result.ExecuteInfo.Job.Command,
			Output:       string(result.Output),
			PlanTime:     result.ExecuteInfo.PlanTime.UnixNano() / 1000 / 1000,
			ScheduleTime: result.ExecuteInfo.RealTime.UnixNano() / 1000 / 1000,
			StartTime:    result.StartTime.UnixNano() / 1000 / 1000,
			EndTime:      result.EndTime.UnixNano() / 1000 / 1000,
		}
		if result.Err != nil {
			jobLog.Err = result.Err.Error()
		} else {
			jobLog.Err = ""
		}
		G_logSink.Append(jobLog)
	}




	//fmt.Println("任务执行完成:",result.ExecuteInfo.Job.Name,string(result.Output),result.Err)
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
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
		jobResultChan: make(chan *common.JobExecuteResult,1000),
	}

	//启动调度协程
	go G_scheduler.schedulerLoop()
	return
}


//回传任务执行结果
func (scheduler *Scheduler) PushJobResult(jobResult *common.JobExecuteResult){
	scheduler.jobResultChan <- jobResult
}