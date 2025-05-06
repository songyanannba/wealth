package timedtask

func TimedTask() {
	// 定时清理任务
	ClearTask()

	// 每日汇总任务
	DailyStatisticsTask()

	// 交互类型的游戏操作
	//InteractiveGame()
}

func ClearTask() {
	// 清理部分过期数据
	//if _, err := global.GVA_Timer.AddTaskByFunc("Clear", "0 0 * * *", Clear); err != nil {
	//	global.GVA_LOG.Error("add timer Clear error:", zap.Error(err))
	//}
	//
	//// 清理保险杠数据 每月1号凌晨0点执行
	//if _, err := global.GVA_Timer.AddTaskByFunc("ClearBumperData", "0 0 1 * *", ClearBumperData); err != nil {
	//	global.GVA_LOG.Error("add timer MoneySlotMonth error:", zap.Error(err))
	//}
	//
	//// 清理在线用户
	//if _, err := global.GVA_Timer.AddTaskByFunc("ClearUserOnline", "30 * * * * *", ClearUserOnlineFn, cron.WithSeconds()); err != nil {
	//	global.GVA_LOG.Error("add timer ClearUserOnlineFn error:", zap.Error(err))
	//}
}

// DailyStatisticsTask 每日汇总
func DailyStatisticsTask() {
	//if _, err := global.GVA_Timer.AddTaskByFunc("DailyStatistics", "5 0 * * *", DailyStatistics); err != nil {
	//	global.GVA_LOG.Error("add timer DailyStatistics error:", zap.Error(err))
	//}
}

// InteractiveGame 交互类型的游戏操作
func InteractiveGame() {
	//if _, err := global.GVA_Timer.AddTaskByFunc("DefaultGamble", "0 */1 * * *", EndInteractiveGame); err != nil {
	//	global.GVA_LOG.Error("add timer DefaultGamble error:", zap.Error(err))
	//}
}
