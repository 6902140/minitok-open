package controllers

import (
	"minitok/response"
	"minitok/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 关注操作
func RelationAction(ctx *gin.Context) {
	//token := ctx.Query("token")
	//tokenUserId, err := usal.VerifyToken(token)
	/* if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	} */
	tokens, _ := ctx.Get("UserId")
	tokenUserId := tokens.(int64)

	toUserId := ctx.Query("to_user_id")
	touid, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	action := ctx.Query("action_type")
	err = services.RelationAction(touid, tokenUserId, action)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", nil)
}

// 获取关注列表
func GetFollowList(ctx *gin.Context) {
	//token := ctx.Query("token")
	//tokenUserId, err := usal.VerifyToken(token)
	/* if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	} */
	tokens, _ := ctx.Get("UserId")
	tokenUserId := tokens.(int64)

	UserId := ctx.Query("user_id")
	uid, err := strconv.ParseInt(UserId, 10, 64)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	followList, err := services.RelationFollowList(uid, tokenUserId)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", followList)
}

// 获取关注者列表
func GetFollowerList(ctx *gin.Context) {
	/* token := ctx.Query("token")
	tokenUserId, err := usal.VerifyToken(token)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	} */
	tokens, _ := ctx.Get("UserId")
	tokenUserId := tokens.(int64)

	UserId := ctx.Query("user_id")
	uid, err := strconv.ParseInt(UserId, 10, 64)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	followerList, err := services.RelationFollowerList(uid, tokenUserId)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", followerList)
}
