package master

import (
	"crontab/common"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"
)

/**
定义任务的HTTP接口
*/
type ApiServer struct {
	httpServer *http.Server
}

/**
定义单列的全局apiserver
*/
var (
	//单例对象
	G_apiServer *ApiServer
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
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		job     common.CronJob
		postJob string
		oldJob  *common.CronJob
		bytes   []byte
	)
	//任务保存在etcd中
	//1.解析POST表单提交
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	//2.去表单中的job字段
	postJob = req.PostForm.Get("job")
	//3.反序列化job,将postJob序列化为字节数组，然后赋值给job
	err = json.Unmarshal([]byte(postJob), &job)
	if err != nil {
		goto ERR
	}
	//4.保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
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
func handleJobDelete(resp http.ResponseWriter,req *http.Request){
	var(
		err error
		name string
		oldJob *common.CronJob
		bytes []byte
	)

	//POST表单数据格式(a=1& b= 2 & c=3
	if err = req.ParseForm(); err != nil{
		goto ERR
	}
	//删除的任务名
	name = req.PostForm.Get("name")

	//删除任务
	if oldJob,err = G_jobMgr.DeleteJob(name); err != nil{
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
func handleJobList(resp http.ResponseWriter,req *http.Request){
	var(
		err error
		jobList []*common.CronJob
		bytes []byte
	)

	//查询任务列表
	if jobList, err = G_jobMgr.ListJob(); err != nil{
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
func handleJobKill(resp http.ResponseWriter, req *http.Request){
	var(
		err error
		name string
		bytes []byte
	)
	//解析表单
	if err = req.ParseForm(); err != nil{
		goto ERR
	}
	//要杀死任务的任务名
	name = req.PostForm.Get("name")

	if err = G_jobMgr.KillJob(name); err != nil{
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
//初始化服务
func InitApiServer() (err error) {
	var(
		staticDir http.Dir
		staticHandler http.Handler
	)
	//配置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/cron/job/save", handleJobSave)
	mux.HandleFunc("/cron/job/delete", handleJobDelete)
	mux.HandleFunc("/cron/job/list",handleJobList)
	mux.HandleFunc("/cron/job/kill",handleJobKill)
	// 首页请求路由: /index.html
	//静态文件目录
	staticDir = http.Dir(G_config.WebRoot)

	staticHandler = http.FileServer(staticDir)
	//匹配到/index.html后会去掉/然后加上./webroot形成./webroot/index.html
	mux.Handle("/",http.StripPrefix("/",staticHandler))
	//启动tcp监听地址和端口
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort))
	if err != nil {
		return
	}

	//创建http服务
	httpServer := &http.Server{
		//定义http读写超时时间
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}
	G_apiServer = &ApiServer{httpServer: httpServer,}

	//让服务启动在协程中
	go httpServer.Serve(listener)
	return
}
