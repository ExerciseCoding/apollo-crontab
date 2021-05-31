package common

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

//应答方法
func BuildResponse