package routers

import (
	"gateway/middleware"
	"github.com/gin-gonic/gin"
)

func MemeInit(router *gin.Engine) {
	//magicRouter := router.Group("/mt_way").Use(middleware.JWTAuth(), middleware.ServiceConfig())
	magicRouter := router.Group("/mm_way").Use(middleware.JWTAuth())
	{

		magicRouter.POST("/betList", func(context *gin.Context) {

		})

	}

}
