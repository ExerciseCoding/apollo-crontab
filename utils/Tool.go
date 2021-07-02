package utils

import (
	"net"

	"apollo/common"
)

//获取本机网卡IP
func GetLocalIP()(ipv4 string, err error){
	var(
		addrs []net.Addr
		addr net.Addr
		ipNet *net.IPNet
		isIPNet bool
	)
	//遍历所有网卡
	if addrs,err = net.InterfaceAddrs(); err != nil{
		return
	}

	//取第一个非localhost的网卡IP
	for _,addr = range addrs{
		//ipv4,ipv6
		//判断这个网络地址是否是IP地址? ipv4, ipv6
		if ipNet,isIPNet = addr.(*net.IPNet); isIPNet && !ipNet.IP.IsLoopback(){
			//跳过IPV6
			if ipNet.IP.To4() != nil{
				ipv4 = ipNet.IP.String() //198.168.1.1
				return
			}
		}
	}
	err = common.ERR_NO_LOCAL_IP_FOUND
	return

}