package router

import (
	"net/http"

	"apollo/master/controller"
)

func InitRouter()(*http.ServeMux){
	mux := http.NewServeMux()
	mux.HandleFunc("/cron/job/save",controller.G_jobController.HandleJobSave)
	mux.HandleFunc("/cron/job/delete", controller.G_jobController.HandleJobDelete)
	mux.HandleFunc("/cron/job/list",controller.G_jobController.HandleJobList)
	mux.HandleFunc("/cron/job/kill",controller.G_jobController.HandleJobKill)
	mux.HandleFunc("/cron/job/update",controller.G_jobController.HandleJobUpdate)
	mux.HandleFunc("/job/log",controller.G_jobLogController.HandleJobLog)
	mux.HandleFunc("/worker/list",controller.G_workController.HandleWorkerList)
	mux.HandleFunc("/cron/user/login",controller.G_loginController.HandleLogin)
	mux.HandleFunc("/cron/user/info",controller.G_loginController.HandleGetLoginInfo)
	mux.HandleFunc("/cron/user/logout",controller.G_loginController.HandleLoginout)

	return mux
}