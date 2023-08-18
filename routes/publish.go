package routes

import (
	"minitok/controllers"
	"minitok/usal"

	"github.com/gin-gonic/gin"
)

func PublishRoutes(r *gin.RouterGroup) {
	publish := r.Group("publish")
	{
		publish.POST("/action/", usal.AuthMiddleware(), controllers.PublishAction)
		publish.GET("/list/", usal.AuthWithOutMiddleware(), controllers.GetPublishList)
	}
}
