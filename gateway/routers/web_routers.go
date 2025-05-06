// Package routers 路由
package routers

import (
	"gateway/global"
	"gateway/helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
)

// Init http 接口路由
func Init(router *gin.Engine) {
	//router.LoadHTMLGlob("views/**/*")

	//执行上架队列
	//go cache.ExecDequeue(cache.Queue_UserSellOre)
	//docs.SwaggerInfo.BasePath = global.GVA_CONFIG.System.RouterPrefix
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	////router.GET(global.GVA_CONFIG.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//global.GVA_LOG.Info("register swagger handler")

	//url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	//短连接 登陆游戏
	//router.POST("/login", rat.Login)
	//router.GET("/test", rat.Test)

	SlotInit(router)

	//执行一些脚本
	//router.POST("/task", rat.Task)

	router.Static("assets", "../assets")
	//获取图片
	router.GET("/imagers/getPicture/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		if filename == "0" {
			randIntName := helper.RandInt(6) + 1
			filename = strconv.Itoa(randIntName) + ".jpg"
			global.GVA_LOG.Infof("imagers/getPicture  randIntName:%v , filename:%v", randIntName, filename)
		}
		//filePath := filepath.Join("/Users/syn/goProjects/zhuoshuo/nb_game_server/gate_way/assets/picture", "1.jpg")
		//filename = strconv.Itoa(helper.RandInt(6)) + ".jpg"
		filePath := filepath.Join("./assets/picture", filename)

		// 尝试读取文件并返回
		//abs, err := filepath.Abs(filePath)
		//fmt.Println(abs, err)
		if _, err := filepath.Abs(filePath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file path"})
			return
		}
		c.File(filePath)
	})

	//// 用户组
	//userRouter := router.Group("/user")
	//{
	//	userRouter.GET("/list", user.List)
	//	userRouter.GET("/online", user.Online)
	//	userRouter.POST("/sendMessage", user.SendMessage)
	//	userRouter.POST("/sendMessageAll", user.SendMessageAll)
	//}

	//// 系统
	//systemRouter := router.Group("/system")
	//{
	//	systemRouter.GET("/state", systems.Status)
	//}

	//// home
	//homeRouter := router.Group("/home")
	//{
	//	homeRouter.GET("/index", home.Index)
	//}

	//gameWayRouter := router.Group("/g_way")
	//{
	//	gameWayRouter.GET("/index", way.Index)
	//
	//	gameWayRouter.GET("/swingRod", way.Fish)

	//	//挥杆 swingRod
	//	//gameWayRouter.GET("/swingRod", way.SwingRod)
	//	//
	//	////挥杆收获
	//	gameWayRouter.POST("/sr_Harvest", way.SRHarvest)

	//	////挥杆确认
	//	//gameWayRouter.POST("/sr_Confirm", way.SRConfirm)
	//}
}
