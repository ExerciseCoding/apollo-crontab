package controller

import (
	"fmt"
	"net/http"

	"apollo/common"
)

type LoginController struct {

}

var (
	G_loginController *LoginController = &LoginController{}
)
func(loginCtr *LoginController) HandleLogin(resp http.ResponseWriter, req *http.Request){
	var(
		bytes []byte
		err error
	)
	fmt.Println(req.URL)
	if bytes, err = common.BuildResponse(0,"success",""); err == nil{
		resp.Write(bytes)
	}
	return
}

func(loginCtr *LoginController) HandleGetLoginInfo(resp http.ResponseWriter, req *http.Request){
	var(
		bytes []byte
		err error
	)
	roles := []string{"admin"}
	userInfo := make(map[string]interface{})
	userInfo["id"] = "1"
	userInfo["name"] = "amdin"
	userInfo["roles"] = roles
	userInfo["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	fmt.Println(req.URL)
	if bytes, err = common.BuildResponse(0,"success",userInfo); err == nil{
		resp.Write(bytes)
	}
	return
}

func(loginCtr *LoginController) HandleLoginout(resp http.ResponseWriter, req *http.Request) {
	var (
		bytes []byte
		err   error
	)
	if bytes, err = common.BuildResponse(0, "success", ""); err == nil {
		resp.Write(bytes)
	}
	return
}