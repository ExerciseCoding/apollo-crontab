package master

import (
	"apollo/master/router"
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

//初始化服务
func InitApiServer() (err error) {
	var(
		staticDir http.Dir
		staticHandler http.Handler
		mux *http.ServeMux
	)
	//配置路由
	mux = router.InitRouter()


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
