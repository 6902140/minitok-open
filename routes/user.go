package routes

import (
	"minitok/controllers"
	"minitok/usal"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	//函数创建一个名为 user 的路由组
	user := r.Group("user")
	{

		user.POST("/login/", controllers.UserLogin)
		user.GET("/", usal.AuthMiddleware(), controllers.GetUserInfo)

		user.POST("/register/", controllers.UserRegister)
	}

}
