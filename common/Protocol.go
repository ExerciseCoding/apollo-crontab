package common

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gorhill/cronexpr"
)

//定时任务
type CronJob struct {
	Name string `json:"name"` //任务名
	Command string `json:"command"` //shell命令
	CronExpr string `json:"cronExpr"` //cron表达式
}

//HTTP接口应答
type Response struct {
	Errno int `json:"errno"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}
//变化事件
type JobEvent struct {
	EventType int
	Job *CronJob
}

//任务调度计划
type JobSchedulerPlan struct {
	Job *CronJob //调度的任务
	Expr *cronexpr.Expression //解析好的cronexpr表达式
	NextTime time.Time
}

//应答方法
func BuildResponse(errno int,msg string,data interface{})(resp []byte,err error){
	//1.定义一个Response
	var(
		response Response
	)
	response.Errno = errno
	response.Msg = msg
	response.Data = data

	//序列化
	resp,err = json.Marshal(response)
	return
}

//反序列化Job
func UnpackJob(jobValue []byte)(job *CronJob,err error){
	job = &CronJob{}
	if err = json.Unmarshal(jobValue,job); err != nil {
		return
	}
	return
}

//从etcd的key中提取任务名
//例：从/cron/jobs/job10 中提取到job10
func ExtractJobName(jobKey string)(string){
	return strings.TrimPrefix(jobKey,JOB_SAVE_DIR)
}


//任务变化事件 1.更新任务 2.删除任务
func BuildJobEvent(eventType int,job *CronJob)(jobEvent *JobEvent){
	return &JobEvent{
		EventType:eventType,
		Job:job,
	}
}


func BuildJobSchedulerPlan(job *CronJob)(jobSchedulerPlan *JobSchedulerPlan,err error){
	var(
		expr *cronexpr.Expression
	)
	if expr,err = cronexpr.Parse(job.CronExpr); err != nil{
		return
	}
	jobSchedulerPlan = &JobSchedulerPlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}

	return
}

