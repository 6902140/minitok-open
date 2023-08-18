package controllers

import (
	"fmt"
	"minitok/config"
	"minitok/log"
	"minitok/response"
	"minitok/services"
	"minitok/util"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 视频发布
func PublishAction(ctx *gin.Context) {
	userId, _ := ctx.Get("UserId") //从之前的函数获取userid

	title := ctx.PostForm("title")    //返回请求中名为"title"的表单字段的值。
	data, err := ctx.FormFile("data") //从HTTP POST请求中获取上传的文件。
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	filename := filepath.Base(data.Filename)

	finalName := fmt.Sprintf("%s_%s", util.RandomString(), filename)
	videoPath := config.GetConfig().Path.Videofile
	saveFile := filepath.Join(videoPath, finalName) //结合生成最终的文件地址

	log.Info("saveFile:", saveFile) //打印日志输出相关信息

	//将上传的文件保存到指定的容器的地址
	if err := ctx.SaveUploadedFile(data, saveFile); err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	publish, err := services.PublishVideo(userId.(int64), saveFile, title)
	//publish, err := services.PublishVideo(userId, saveFile)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	log.Infof("publish:%v", publish)
	response.Success(ctx, "success", publish)

}

// 获取视频列表
func GetPublishList(ctx *gin.Context) {
	tokenUserId, _ := ctx.Get("UserId")
	id := ctx.Query("user_id")
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
	}
	list, err := services.PublishList(tokenUserId.(int64), userId)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", list)
}
