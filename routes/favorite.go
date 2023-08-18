package routes

import (
	"minitok/controllers"
	"minitok/usal"

	"github.com/gin-gonic/gin"
)

func FavoriteRoutes(r *gin.RouterGroup) {
	favorite := r.Group("favorite")
	{
		favorite.POST("/action/", usal.AuthMiddleware(), controllers.FavoriteAction)
		favorite.GET("/list/", usal.AuthWithOutMiddleware(), controllers.GetFavoriteList)
	}

}
