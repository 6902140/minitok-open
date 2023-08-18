package routes

import (
	"minitok/controllers"
	"minitok/usal"

	"github.com/gin-gonic/gin"
)

func RelationRoutes(r *gin.RouterGroup) {
	relation := r.Group("relation")
	{
		relation.POST("/action/", usal.AuthMiddleware(), controllers.RelationAction)
		relation.GET("/follow/list/", usal.AuthWithOutMiddleware(), controllers.GetFollowList)
		relation.GET("/follower/list/", usal.AuthWithOutMiddleware(), controllers.GetFollowerList)
	}
}
