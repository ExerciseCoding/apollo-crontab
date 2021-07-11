package controller

import (
	"apollo/common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"apollo/master"
)

type JobController struct {

}
var (
	G_jobController *JobController = &JobController{}
)

//保存任务
/**
POST
job = {
"name":"job1",
"command":"ls -l /root",
"cronExpr": "* * * * *"
}
*/
func (jobCtr *JobController) HandleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		job     common.CronJob
		//postJob string
		oldJob  *common.CronJob
		bytes   []byte
		body     []byte
	)
	if body,err  = ioutil.ReadAll(req.Body); err != nil{
		goto ERR
	}

	fmt.Println(fmt.Sprintf("url"+"%s",req.URL))
	//任务保存在etcd中
	//1.解析POST表单提交
	//if err = req.ParseForm();  err != nil {
	//	goto ERR
	//}

	//2.去表单中的job字段
	//postJob = req.PostForm.Get("job")
	//3.反序列化job,将postJob序列化为字节数组，然后赋值给job
	//err = json.Unmarshal([]byte(postJob), &job)
	err = json.Unmarshal(body, &job)
	if err != nil {
		goto ERR
	}
	//4.保存到etcd
	if oldJob, err = master.G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	//5.返回正常应答({"error":0,"msg":"","data":{......}})
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)

	}
	return

ERR:
	//返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)

	}

}



func(jobCtr *JobController) HandleJobUpdate(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		job     common.CronJob
		oldJob  *common.CronJob
		bytes   []byte
		body     []byte
	)
	if body,err  = ioutil.ReadAll(req.Body); err != nil{
		goto ERR
	}
	fmt.Println(fmt.Sprintf("url"+"%s",req.URL))
	err = json.Unmarshal(body, &job)
	if err != nil {
		goto ERR
	}
	//4.保存到etcd
	if oldJob, err = master.G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	//5.返回正常应答({"error":0,"msg":"","data":{......}})
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)

	}
	return

ERR:
	//返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)

	}

}


//删除任务
//POST /job/delete name=job1
func(jobCtr *JobController) HandleJobDelete(resp http.ResponseWriter,req *http.Request){
	var(
		err error
		name string
		oldJob *common.CronJob
		bytes []byte
		body []byte
		postData map[string]string
		//ok bool
	)
	if body,err  = ioutil.ReadAll(req.Body); err != nil{
		goto ERR
	}
	//name = string(body)
	fmt.Println(req.URL)

	if err = json.Unmarshal(body, &postData); err != nil{
		goto ERR
	}
	name = postData["name"]
	//POST表单数据格式(a=1& b= 2 & c=3
	//if err = req.ParseForm(); err != nil{
	//	goto ERR
	//}
	//删除的任务名
	//name = req.PostForm.Get("name")

	//删除任务
	if oldJob,err = master.G_jobMgr.DeleteJob(name); err != nil{
		goto  ERR
	}

	if bytes,err = common.BuildResponse(0,"success",oldJob); err == nil{
		resp.Write(bytes)
	}
	return

ERR:
	if bytes,err = common.BuildResponse(-1,err.Error(),nil); err == nil{
		resp.Write(bytes)
	}
}
//查看所有任务
func(jobCtr *JobController) HandleJobList(resp http.ResponseWriter,req *http.Request){
	var(
		err error
		jobList []*common.CronJob
		bytes []byte
	)

	//查询任务列表
	if jobList, err = master.G_jobMgr.ListJob(); err != nil{
		goto ERR
	}
	if bytes, err = common.BuildResponse(0,"success",jobList); err == nil{
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1,err.Error(),nil); err == nil{
		resp.Write(bytes)
	}

}
//杀死任务
//POST /job/kill name=job1
func(jobCtr *JobController) HandleJobKill(resp http.ResponseWriter, req *http.Request){
	var(
		err error
		name string
		bytes []byte
		body []byte
		postData map[string]string
	)
	if body,err  = ioutil.ReadAll(req.Body); err != nil{
		goto ERR
	}
	//name = string(body)
	fmt.Println(req.URL)

	if err = json.Unmarshal(body, &postData); err != nil{
		goto ERR
	}
	name = postData["name"]
	fmt.Println(name)

	//解析表单
	//if err = req.ParseForm(); err != nil{
	//	goto ERR
	//}
	//要杀死任务的任务名
	//name = req.PostForm.Get("name")

	if err = master.G_jobMgr.KillJob(name); err != nil{
		goto ERR
	}
	if bytes, err = common.BuildResponse(0,"success",nil); err == nil{
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1,err.Error(),nil); err == nil{
		resp.Write(bytes)
	}

}
