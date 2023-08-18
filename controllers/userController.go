package controllers

import (
	"minitok/log"
	"minitok/response"
	"minitok/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserLogin 用户登录
func UserLogin(context *gin.Context) {
	var err error
	userName := context.Query("username")
	password := context.Query("password")
	//从接口获得用户名称和密码
	if len(userName) > 24 || len(password) > 24 {
		response.Fail(context, "username or password invalid", nil)
		return
	}
	loginResponse, err := services.UserLogin(userName, password) //调用service包中的login函数获取基本信息
	if err != nil {
		log.Infof("login error : %s", err)
		response.Fail(context, err.Error(), nil)
		return
	}
	response.Success(context, "success", loginResponse)
}

// UserRegister 用户注册模块函数
func UserRegister(context *gin.Context) {
	var err error
	userName := context.Query("username")
	password := context.Query("password")
	if len(userName) > 32 || len(password) > 32 { //最长32位字符
		response.Fail(context, "username or password invalid", nil) //用户名或者密码长度不合规则直接返回错误信息即可
		return
	}
	//调用services层的用户注册函数
	registResponse, err := services.UserRegister(userName, password)
	if err != nil {
		log.Infof("registe error : %s", err)
		response.Fail(context, err.Error(), nil)
		return
	}
	response.Success(context, "success", registResponse)

}

// GetUserInfo 获取用户信息
func GetUserInfo(context *gin.Context) {

	var err error
	userId := context.Query("user_id")
	uids, _ := context.Get("UserId")

	uid := uids.(int64)
	if err != nil {
		response.Fail(context, err.Error(), nil)
		return
	}

	if strconv.FormatInt(uid, 10) != userId {
		response.Fail(context, "token error", nil)
		return
	}
	userinfo, err := services.UserInfo(uid)
	if err != nil {
		log.Infof("get userinfo  error : %s", err)
		response.Fail(context, err.Error(), nil)
		return
	}
	response.Success(context, "success", userinfo)

}
