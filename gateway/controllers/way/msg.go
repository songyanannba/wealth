package way

//var battleRoyaleService = logic.OneGetawayBetValue{}
//
//func Fish(c *gin.Context) {
//	data := gin.H{
//		"title": "DealBetON",
//	}
//	global.GVA_LOG.Infof("SwingRod", data)
//	//logic.DealExceptionBetOn()
//
//	controllers.Response(c, common.OK, "", data)
//}
//
//// SwingRod ==
//func SwingRod(c *gin.Context) {
//	data := gin.H{
//		"title": "DealBetON",
//	}
//	global.GVA_LOG.Infof("SwingRod", data)
//	//logic.DealExceptionBetOn()
//
//	controllers.Response(c, common.OK, "", data)
//}
//
//func SRHarvest(c *gin.Context) {
//	data := gin.H{
//		"title": "DealBetON",
//	}
//	controllers.Response(c, common.OK, "", data)
//}
//
//func SRConfirm(c *gin.Context) {
//	data := gin.H{
//		"title": "DealBetON",
//	}
//	controllers.Response(c, common.OK, "", data)
//}
//
//func Index(c *gin.Context) {
//	appIDStr := c.Query("appID")
//	appIDUint64, _ := strconv.ParseInt(appIDStr, 10, 32)
//	appID := uint32(appIDUint64)
//	if !websocket.InAppIDs(appID) {
//		appID = websocket.GetDefaultAppID()
//	}
//	fmt.Println("http_request 聊天首页", appID)
//	data := gin.H{
//		"title":        "聊天首页",
//		"appID":        appID,
//		"httpUrl":      viper.GetString("app.httpUrl"),
//		"webSocketUrl": viper.GetString("app.webSocketUrl"),
//	}
//
//	global.GVA_LOG.Infof("Index===", data)
//
//	//data := make(map[string]interface{})
//
//	controllers.Response(c, common.OK, "", data)
//}
//
//func GetConfig(c *gin.Context) {
//	//获取参数
//
//	//数据库获取数据
//	//res, err := battleRoyaleService.GetawayBetValue()
//	//if err != nil {
//	//	global.GVA_LOG.Infof("GetConfig : ", zap.Error(err))
//	//	return
//	//}
//	//
//	////返回数据
//	//var exts []models.BetExt
//	//
//	//err = json.Unmarshal([]byte(res.Ext), &exts)
//	//if err != nil {
//	//	global.GVA_LOG.Infof("GetConfig Unmarshal ", zap.Error(err))
//	//	return
//	//}
//	//
//	//var extVal []int
//	//for _, val := range exts {
//	//	v, _ := strconv.Atoi(val.Val)
//	//	extVal = append(extVal, v)
//	//}
//
//	data := make(map[string]interface{})
//	//data["icon"] = res.Icon
//	//data["bet"] = extVal
//	//data["wait_countdown"] = res.WaitCountdown //选择房间倒计时
//	////data["next_play_time"] = res.NextPlayTime
//	//data["winner_take_percentage"] = res.WinnerTakePercentage //胜者抽成比例
//
//	controllers.Response(c, common.OK, "", data)
//}
