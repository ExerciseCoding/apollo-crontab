package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"apollo/worker"
)

var(
	confFile string //配置文件路径
)
//初始化线程
func initEnv(){
	//配置程序的最大线程数和CPU的核心数相同
	runtime.GOMAXPROCS(runtime.NumCPU())
}

//解析命令行参数
func initArgs(){
	//worker -config ./worker.json
	//worker -h
	flag.StringVar(&confFile,"config","./worker.json","指定worker.json")
	flag.Parse()
}
func main(){
	var (
		err error
	)
	//初始化命令行参数
	initArgs()
	//初始化线程
	initEnv()
	//加载配置
	if err = worker.InitConfig(confFile); err != nil{
		goto ERR
	}

	if err = worker.InitRegister(); err != nil{
		goto ERR
	}
	//启动日志协程
	if err = worker.InitLogSink(); err != nil{
		goto ERR
	}
	//启动执行器
	if err = worker.InitExecutor(); err != nil {
		goto ERR
	}
	//启动调度器
	if err = worker.InitScheduler();err != nil{
		goto ERR
	}
	//初始化任务管理器
	if err = worker.InitJobMgr(); err != nil{
		goto ERR
	}


	for {
		time.Sleep(1 * time.Second)
	}

	//正常退出
	return
ERR:
	fmt.Println(err)

}