package controllers

import (
	"minitok/response"
	"minitok/services"
	"minitok/util"
	"strconv"

	// "encoding/json"
	"github.com/gin-gonic/gin"
)

// 视频流
func Feed(ctx *gin.Context) {
	var userId int64
	currentTime, err := strconv.ParseInt(ctx.Query("latest_time"), 10, 64)
	if err != nil || currentTime == int64(0) {
		currentTime = util.GetCurrentTime()
	}
	//token := ctx.Query("token")
	//userId, err = usal.VerifyToken(token)
	userIds, _ := ctx.Get("UserId")
	userId = userIds.(int64)

	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	//调用service层的GetFeedList函数获取视频推送信息
	feedList, err := services.GetFeedList(currentTime, userId)
	if err != nil { //获取视频失败时返回信息
		response.Fail(ctx, err.Error(), nil)
	}
	response.Success(ctx, "success", feedList)
}
