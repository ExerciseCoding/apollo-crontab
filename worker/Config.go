package worker

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	EtcdEndpoints []string `json:"etcdEndpoints"`
	EtcdDialTimeout int `json:"etcdDialTimeout"`
}
var (
	G_config *Config
)
//加载服务配置
func InitConfig(filename string)(err error){
	var(
		config Config
	)
	//把配置文件读进来
	context , err := ioutil.ReadFile(filename)
	if err != nil{
		return
	}

	//2.把从配置文件读出来的二进制数组转成json(反序列化)
	if err = json.Unmarshal(context,&config); err != nil{
		return
	}
	//赋值单例
	G_config = &config
	return
}