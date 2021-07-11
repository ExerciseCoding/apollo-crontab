package controller

import (
	"apollo/common"
	"apollo/master"
	"fmt"
	"net/http"
	"strconv"
)

type JobLogController struct {

}

var (
	G_jobLogController *JobLogController = &JobLogController{}
)

//查看日志
func(jobLogCtr *JobLogController) HandleJobLog(resp http.ResponseWriter, req *http.Request){
	var(
		err error
		name string //任务名
		skipParam string//从第几条开始
		limitParam string//返回多少条
		skip int
		limit int
		logArr []*common.JobLog
		bytes []byte
	)
	//解析GET参数
	if err = req.ParseForm(); err != nil{
		goto ERR
	}
	//获取请求参数 /job/log?name=job10&skip=0&limit=10
	name = req.Form.Get("name")
	skipParam = req.Form.Get("skip")
	limitParam = req.Form.Get("limit")

	if skip,err = strconv.Atoi(skipParam); err != nil{
		skip = 0
	}

	if limit, err = strconv.Atoi(limitParam); err != nil{
		limit = 20
	}
	fmt.Println(name,skip,limit)
	if logArr,err = master.G_logMgr.ListLog(name,skip,limit); err != nil{
		goto  ERR
	}
	fmt.Println(logArr)
	if bytes, err = common.BuildResponse(0,"success",logArr); err == nil{
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1,err.Error(),nil); err == nil{
		resp.Write(bytes)
	}

}


