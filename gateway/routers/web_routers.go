// Package routers 路由
package routers

import (
	"github.com/gin-gonic/gin"
)

// Init http 接口路由
func Init(router *gin.Engine) {
	SlotInit(router)
}
