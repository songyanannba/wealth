package core

import (
	"fmt"
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/initialize"

	"time"
)

func BaseInit() {
	ConfInit()

	//mysql
	//global.GVA_USER_DB = initialize.Gorm() // gorm连接数据库

	global.GVA_SLOT_SERVER_DB = initialize.Gorm() // gorm连接数据库

	//定时任务
	initialize.Timer()

	// 初始化redis服务
	initialize.Redis()
}

func ConfInit() {
	global.GVA_VP = Viper() // 初始化Viper
	global.GVA_LOG = Zap()  // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG.Logger)
	global.Location, _ = time.LoadLocation("Asia/Shanghai")
}

func CloseDB() {
	db, _ := global.GVA_SLOT_SERVER_DB.DB()
	db.Close()
	fmt.Println("关闭读写数据库连接")
	//CloseTavernStoryDB()
}

//func CloseTavernStoryDB() {
//	TSdb, _ := global.GVA_SLOT_SERVER_DB.DB()
//	TSdb.Close()
//	global.GVA_LOG.Infof("GVA_SLOT_SERVER_DB 关闭数据库连接")
//
//}
