package core

import (
	"gateway/global"
	"gateway/initialize"
	"gateway/utils/queue"
	"go.uber.org/zap"
)

func BaseInit() {
	ConfInit()
	//mysql
	global.GVA_USER_DB = initialize.Gorm() // gorm连接数据库
	//global.GVA_MEME_DB = initialize.MemeBattleMysql() // gorm连接数据库
	//initialize.DBList()

	initialize.Timer()

	// 初始化redis服务
	initialize.Redis()

	// 创建一个 的队列
	global.ChanQueue = queue.NewQueue(100000)
	global.QueueDataKeyMap = queue.NewQueueDataKeyMap()
}

func ConfInit() {
	global.GVA_VP = Viper() // 初始化Viper
	global.GVA_LOG = Zap()  // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG.Logger)
}

func CloseDB() {
	//
	db, _ := global.GVA_USER_DB.DB()
	db.Close()
	global.GVA_LOG.Infof("GVA_USER_DB 关闭数据库连接")

	//骗子酒馆库
	//CloseMemeDB()
}

//func CloseMemeDB() {
//
//	TSdb, _ := global.GVA_MEME_DB.DB()
//	TSdb.Close()
//	global.GVA_LOG.Infof("GVA_MEME_DB 关闭数据库连接")
//
//}
