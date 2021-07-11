package controller

import (
	"apollo/common"
	"apollo/master"
	"net/http"
)

type WorkController struct {

}

var(
	G_workController *WorkController = &WorkController{}
)
//获取监控worker节点列表
func(workCtr *WorkController) HandleWorkerList(resp http.ResponseWriter, req *http.Request){
	var(
		workerArr []string
		bytes []byte
		err error
	)
	if workerArr, err = master.G_workerMgr.ListWorkers(); err != nil{
		goto ERR
	}

	if bytes, err = common.BuildResponse(0,"success",workerArr); err == nil{
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1,err.Error(),nil); err == nil{
		resp.Write(bytes)
	}

}