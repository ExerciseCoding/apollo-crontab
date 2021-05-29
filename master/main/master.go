package main

import "runtime"

func initEnv(){
	//配置程序的最大线程数和CPU的核心数相同
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main(){
	//初始化线程
	initEnv()

	//启动Api http服务
}