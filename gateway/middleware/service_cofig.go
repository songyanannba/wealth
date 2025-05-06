package middleware

import (
	"gateway/common"
	"gateway/controllers"
	"gateway/global"
	"gateway/servers/src/dao"
	"gateway/utils/cache"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ServiceConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取标识 是否在服务中
		//服务维护阶段
		val, err := cache.GetGetServiceConfigKeyExpPre()
		if err != nil {
			global.GVA_LOG.Error("GetGetServiceConfigKeyExpPre", zap.Any("err", err))
		}
		if len(val) > 0 {
			if val == "1" {
				controllers.Response(c, common.Maintenance, "", nil)
				c.Abort()
				return
			}
		}

		gameServiceConf, err := dao.GetGameServiceConf(2)
		if err != nil {
			global.GVA_LOG.Error("ServiceConfig", zap.Error(err))
			controllers.Response(c, common.ServerError, "", nil)
			c.Abort()
			return
		}

		if gameServiceConf.Maintenance == 1 {
			global.GVA_LOG.Infof("ServiceConfig %v", *gameServiceConf)
			controllers.Response(c, common.Maintenance, "", nil)
			c.Abort()
			return
		}

		err = cache.SetGetServiceConfigKeyExpPre(gameServiceConf.Maintenance)
		if err != nil {
			global.GVA_LOG.Error("ServiceConfig", zap.Error(err))
		}

		c.Next()
	}
}
