package master

import (
	"crontab/common"
	"encoding/json"
	"fmt"
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
var(
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
func handleJobSave(rw http.ResponseWriter, req *http.Request){
	var(
		err error
		job common.CronJob
		postJob string
	)
	//任务保存在etcd中
	//1.解析POST表单提交
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	//2.去表单中的job字段
	postJob = req.PostForm.Get("job")
	err = json.Unmarshal([]byte(postJob),&job)
	if err != nil{
		goto ERR
	}
	//反序列化job,将postJob序列化为字节数组，然后赋值给job


ERR:
	fmt.Println(err)

}
//初始化服务
func InitApiServer()(err error){
	//配置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/cron/job/save",handleJobSave)

	//启动tcp监听地址和端口
	 listener,err := net.Listen("tcp",":"+strconv.Itoa(G_config.ApiPort))
	 if err != nil{
		return
	}

	//创建http服务
	httpServer := &http.Server{
		//定义http读写超时时间
		ReadTimeout: time.Duration(G_config.ApiReadTimeout)*time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout)*time.Millisecond,
		Handler:mux,
	}
	G_apiServer = &ApiServer{httpServer:httpServer,}

	//让服务启动在协程中
	go httpServer.Serve(listener)
	return
}