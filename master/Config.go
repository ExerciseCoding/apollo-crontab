package master

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ApiPort int `json:"apiPort"`
	ApiReadTimeout int `json:"apiReadTimeOut"`
	ApiWriteTimeout int `json:"apiWriteTimeout"`
	EtcdEndpoints []string `json:"etcdEndpoints"`
	EtcdDialTimeout int `json:"etcdDialTimeout"`
	WebRoot string `json:"webroot"`
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