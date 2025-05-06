package routers

import (
	"gateway/controllers/user"
	"gateway/middleware"
	"github.com/gin-gonic/gin"
)

func SlotInit(router *gin.Engine) {

	//注册
	router.POST("/register", user.Register)
	//登录
	router.POST("/login", user.Login)

	//magicRouter := router.Group("/mt_way").Use(middleware.JWTAuth(), middleware.ServiceConfig())
	magicRouter := router.Group("/game_way").Use(middleware.JWTAuth())
	{

		magicRouter.GET("/user_info", user.GUserInfo)

	}

}
