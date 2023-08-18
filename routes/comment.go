package routes

import (
	"minitok/controllers"
	"minitok/usal"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(r *gin.RouterGroup) {
	comment := r.Group("comment")
	{
		comment.POST("/action/", usal.AuthMiddleware(), controllers.CommentAction)
		comment.GET("/list/", usal.AuthWithOutMiddleware(), controllers.GetCommentList)
	}

}
