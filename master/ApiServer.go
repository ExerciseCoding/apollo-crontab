package master

import (
	"net"
	"net/http"
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
func handleJobSave(rw http.ResponseWriter, req *http.Request){

}
//初始化服务
func initApiServer()(err error){
	//配置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/cron/job/save",handleJobSave)

	//启动tcp监听地址和端口
	 listener,err := net.Listen("tcp",":9999")
	 if err != nil{
		return
	}

	//创建http服务
	httpServer := &http.Server{
		//定义http读写超时时间
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:mux,
	}
	G_apiServer = &ApiServer{httpServer:httpServer,}

	//让服务启动在协程中
	go httpServer.Serve(listener)
	return
}