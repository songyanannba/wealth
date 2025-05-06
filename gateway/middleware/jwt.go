package middleware

import (
	"gateway/common"
	"gateway/controllers"
	"gateway/global"
	"gateway/helper"
	"github.com/gin-gonic/gin"
)

//func SimpleMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 在方法执行前执行的操作
//		println("Before request")
//
//		// 继续执行其他的链路
//		c.Next()
//
//		// 在方法执行后执行的操作
//		println("After request")
//	}
//}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//我们这里jwt鉴权取头部信息 gw-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("gw-token")
		if token == "" {
			//非法的用户
			controllers.Response(c, common.Unauthorized, "", nil)
			c.Abort()
			return
		}

		//jwt 验证
		claims, err := helper.ParseJWT(token)
		global.GVA_LOG.Infof("ParseJWT err %v ", claims)
		if err != nil {
			//未授权
			global.GVA_LOG.Infof(" ParseJWT err ")
			controllers.Response(c, common.UnauthorizedUserToken, "", nil)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
