package meme_battle

//
//// MtBetList 魔塔手动押注列表
//// @Summary      魔塔手动押注列表
//// @Tags         魔塔
//// @Description  魔塔手动押注列表
//// @Accept       json
//// @Produce      json
//// @Param     name query string true "用户姓名"
//// @Param        user  body      models.MtBetReq                true  "魔塔手动押注列表"
//// @Success  1   {object}        common.JSONResult{data=models.MtBetListResp} "魔塔手动押注列表"
//// @Failure      400   {object}  common.JSONResult                "错误提示"
//// @Router       /betList [post]
//func MtBetList(c *gin.Context) {
//	var (
//		request = &models.MtBetReq{}
//		data    = make(map[string]interface{})
//		//resp = models.MtBetListResp{}
//	)
//
//	err := c.ShouldBindJSON(&request)
//	if err != nil {
//		controllers.Response(c, common.ParameterIllegal, "", data)
//		return
//	}
//
//	userId, err := helper.GetUserID(c)
//	if err != nil {
//		controllers.Response(c, common.UnauthorizedUserID, "", data)
//		return
//	}
//
//	reqData := &pbs.BetReq{
//		UserId: userId,
//	}
//	reqMarshal, _ := proto.Marshal(reqData)
//
//	idempotent, err := cache.GetIdempotent(userId, "MtBetList")
//	if err != nil {
//		controllers.Response(c, common.ServerError, "", data)
//		return
//	}
//	if len(idempotent) > 0 {
//		controllers.Response(c, common.DuplicateRequests, "", data)
//		return
//	}
//	err = cache.SetIdempotentNx(userId, "MtBetList", "MtBetList")
//	if err != nil {
//		controllers.Response(c, common.ServerError, "", data)
//		return
//	}
//
//	// 调用 gRPC 方法
//	msgReq := pbs.NetMessage{
//		ReqHead: &pbs.ReqHead{
//			Uid:      0,
//			Token:    "",
//			Platform: "",
//		},
//		AckHead:   &pbs.AckHead{},
//		ServiceId: "",
//		MsgId:     int32(pbs.Meb_pnBetReq),
//		Content:   reqMarshal,
//	}
//	response, err := grpcclient.GetMtClient().CallMtMethod(&msgReq)
//	if response != nil && response.AckHead.Code != pbs.Code_OK {
//		controllers.Response(c, uint32(response.AckHead.Code), "", data)
//		return
//	}
//	if err != nil || response == nil {
//		global.GVA_LOG.Error("could not call method:", zap.Error(err))
//		controllers.Response(c, common.RpcCallRespErr, "", data)
//		return
//	}
//
//	var betListAck pbs.BetListAck
//	respData := response.Content
//	err = proto.Unmarshal(respData, &betListAck)
//	if err != nil {
//		global.GVA_LOG.Error("Unmarshal BetListData :", zap.Error(err))
//		controllers.Response(c, common.RpcCallRespDataErr, "", data)
//		return
//	}
//
//	global.GVA_LOG.Infof("betListAck: %v", &betListAck)
//
//	var betList []models.BetList
//
//	for k, _ := range betListAck.BetListData {
//		listData := betListAck.BetListData[k]
//		item := models.BetList{
//			Id:  int(listData.Id),
//			Bet: float64(listData.Bet),
//		}
//		betList = append(betList, item)
//	}
//
//	data = gin.H{
//		"bet_list": betList,
//	}
//
//	global.GVA_LOG.Infof("MtBetList data:%v", data)
//	controllers.Response(c, common.WebOK, "", data)
//}
//
//// MtCurrStatus  魔塔当前状态
//// @Summary      魔塔当前状态
//// @Tags         魔塔
//// @Description  魔塔当前状态
//// @Accept       json
//// @Produce      json
//// @Param        name query string true "用户姓名"
//// @Param        user  body      models.MTStatusReq                true  "魔塔当前状态"
//// @Success  1   {object}        common.JSONResult{data=models.MTStatusResp} "魔塔当前状态"
//// @Failure      400   {object}  common.JSONResult                "错误提示"
//// @Router       /mtCurrStatus [post]
//func MtCurrStatus(c *gin.Context) {
//	var (
//		request = &models.MTStatusReq{}
//		data    = make(map[string]interface{})
//	)
//
//	err := c.ShouldBindJSON(&request)
//	if err != nil {
//		controllers.Response(c, common.ParameterIllegal, "", data)
//		return
//	}
//
//	userId, err := helper.GetUserID(c)
//	if err != nil {
//		controllers.Response(c, common.UnauthorizedUserID, "", data)
//		return
//	}
//
//	reqData := &pbs.MTStatusReq{
//		UserId: userId,
//	}
//	reqMarshal, _ := proto.Marshal(reqData)
//
//	idempotent, err := cache.GetIdempotent(userId, "MtCurrStatus")
//	if err != nil {
//		controllers.Response(c, common.ServerError, "", data)
//		return
//	}
//	if len(idempotent) > 0 {
//		controllers.Response(c, common.DuplicateRequests, "", data)
//		return
//	}
//	err = cache.SetIdempotentNx(userId, "MtCurrStatus", "MtCurrStatus")
//	if err != nil {
//		controllers.Response(c, common.ServerError, "", data)
//		return
//	}
//
//	// 调用 gRPC 方法
//	msgReq := pbs.NetMessage{
//		ReqHead: &pbs.ReqHead{
//			Uid:      0,
//			Token:    "",
//			Platform: "",
//		},
//		AckHead:   &pbs.AckHead{},
//		ServiceId: "",
//		MsgId:     int32(pbs.Meb_pnStatusReq),
//		Content:   reqMarshal,
//	}
//	response, err := grpcclient.GetMtClient().CallMtMethod(&msgReq)
//	if response != nil && response.AckHead.Code != pbs.Code_OK {
//		controllers.Response(c, uint32(response.AckHead.Code), "", data)
//		return
//	}
//	if err != nil || response == nil {
//		global.GVA_LOG.Error("could not call method:", zap.Error(err))
//		controllers.Response(c, common.RpcCallRespErr, "", data)
//		return
//	}
//
//	var mTStatusAck pbs.MTStatusAck
//	respData := response.Content
//	err = proto.Unmarshal(respData, &mTStatusAck)
//	if err != nil {
//		global.GVA_LOG.Error("Unmarshal mTStatusAck :", zap.Error(err))
//		controllers.Response(c, common.RpcCallRespDataErr, "", data)
//		return
//	}
//
//	global.GVA_LOG.Infof("mTStatusAck: %v", &mTStatusAck)
//
//	resp := models.MTStatusResp{
//		PeriodId:  int(mTStatusAck.PeriodId),
//		State:     int(mTStatusAck.State),
//		Layer:     int(mTStatusAck.Layer),
//		StartTime: mTStatusAck.StartTime,
//		PlayerMeta: models.PlayerMeta{
//			UserId: mTStatusAck.PlayerMeta.UserId,
//		},
//		LayerMeta: models.LayerMeta{},
//	}
//
//	data = gin.H{
//		"period_id":   resp.PeriodId,
//		"state":       resp.State,
//		"layer":       resp.Layer,
//		"start_time":  resp.StartTime,
//		"player_meta": resp.PlayerMeta,
//		"layer_meta":  resp.LayerMeta,
//	}
//
//	global.GVA_LOG.Infof("MtBetList data:%v", data)
//	controllers.Response(c, common.WebOK, "", data)
//}
//
//// MtIsAutoUser  是否成为自动用户
//// @Summary      是否成为自动用户
//// @Tags         魔塔
//// @Description  是否成为自动用户
//// @Accept       json
//// @Produce      json
//// @Param        user  body      models.MTIsAutoReq                true  "是否成为自动用户"
//// @Success  1   {object}        common.JSONResult{data=models.MTIsAutoResp} "是否成为自动用户"
//// @Failure      400   {object}  common.JSONResult                "错误提示"
//// @Router       /isAutoUser [post]
//func MtIsAutoUser(c *gin.Context) {
//	var (
//		request = &models.MTIsAutoReq{}
//		data    = make(map[string]interface{})
//	)
//
//	err := c.ShouldBindJSON(&request)
//	if err != nil {
//		controllers.Response(c, common.ParameterIllegal, "", data)
//		return
//	}
//
//	userId, err := helper.GetUserID(c)
//	if err != nil {
//		controllers.Response(c, common.UnauthorizedUserID, "", data)
//		return
//	}
//
//	reqData := &pbs.IsAutoReq{
//		UserId: userId,
//		IsAuto: int32(request.IsAuto),
//	}
//	reqMarshal, _ := proto.Marshal(reqData)
//
//	msgReq := pbs.NetMessage{
//		ReqHead: &pbs.ReqHead{
//			Uid:      0,
//			Token:    "",
//			Platform: "",
//		},
//		AckHead:   &pbs.AckHead{},
//		ServiceId: "",
//		MsgId:     int32(pbs.Meb_pnIsAutoReq),
//		Content:   reqMarshal,
//	}
//	response, err := grpcclient.GetMtClient().CallMtMethod(&msgReq)
//	if response != nil && response.AckHead.Code != pbs.Code_OK {
//		controllers.Response(c, uint32(response.AckHead.Code), "", data)
//		return
//	}
//	if err != nil || response == nil {
//		global.GVA_LOG.Error("could not call method:", zap.Error(err))
//		controllers.Response(c, common.RpcCallRespErr, "", data)
//		return
//	}
//
//	data = gin.H{}
//	global.GVA_LOG.Infof("MtIsAutoUser data:%v", data)
//	controllers.Response(c, common.WebOK, "", data)
//
//}
//
//// MtBetOnAutoNum 自动用户:充值
//// @Summary       自动用户:充值
//// @Tags         魔塔
//// @Description  自动用户:充值
//// @Accept       json
//// @Produce      json
//// @Param        user  body      models.MTBetOnAutoNumReq                true  "自动用户:充值"
//// @Success  1   {object}        common.JSONResult{data=models.MTBetOnAutoNumResp} "自动用户:充值"
//// @Failure      400   {object}  common.JSONResult                "错误提示"
//// @Router       /betOnAutoNum [post]
//func MtBetOnAutoNum(c *gin.Context) {
//	var (
//		request = &models.MTBetOnAutoNumReq{}
//		data    = make(map[string]interface{})
//	)
//
//	err := c.ShouldBindJSON(&request)
//	if err != nil {
//		controllers.Response(c, common.ParameterIllegal, "", data)
//		return
//	}
//
//	userId, err := helper.GetUserID(c)
//	if err != nil {
//		controllers.Response(c, common.UnauthorizedUserID, "", data)
//		return
//	}
//
//	//本地锁
//	cacheKey := userId + "AutoNum"
//	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
//	if err != nil {
//		global.GVA_LOG.Infof("自动用户 || 一键召唤 QueueDataKeyMap TryAdd%v", cacheKey)
//		controllers.Response(c, common.DuplicateRequests, "", data)
//		return
//	}
//	defer global.QueueDataKeyMap.Del(cacheKey)
//
//	if request.Bet <= 0 {
//		controllers.Response(c, common.ParameterNot, "", data)
//		return
//	}
//
//	reqData := &pbs.AutoNumReq{
//		UserId: userId,
//		Bet:    float32(request.Bet),
//	}
//	reqMarshal, _ := proto.Marshal(reqData)
//
//	msgReq := pbs.NetMessage{
//		ReqHead: &pbs.ReqHead{
//			Uid:      0,
//			Token:    "",
//			Platform: "",
//		},
//		AckHead:   &pbs.AckHead{},
//		ServiceId: "",
//		MsgId:     int32(pbs.Mmb_pnAutoNumReq),
//		Content:   reqMarshal,
//	}
//	response, err := grpcclient.GetMtClient().CallMtMethod(&msgReq)
//	if response != nil && response.AckHead.Code != pbs.Code_OK {
//		controllers.Response(c, uint32(response.AckHead.Code), "", data)
//		return
//	}
//	if err != nil || response == nil {
//		global.GVA_LOG.Error("could not call method:", zap.Error(err))
//		controllers.Response(c, common.RpcCallRespErr, "", data)
//		return
//	}
//
//	data = gin.H{}
//	global.GVA_LOG.Infof("MtBetOnAutoNum data:%v", data)
//	controllers.Response(c, common.WebOK, "", data)
//}
//
//// MtBetOnNum    手动押注：充值
//// @Summary      手动押注：充值
//// @Tags         魔塔
//// @Description 手动押注：充值
//// @Accept       json
//// @Produce      json
//// @Param        user  body      models.MTBetOnNumReq                true  "手动押注：充值"
//// @Success  1   {object}        common.JSONResult{data=models.MTBetOnNumResp} "手动押注：充值"
//// @Failure      400   {object}  common.JSONResult                "错误提示"
//// @Router       /betOnNum [post]
//func MtBetOnNum(c *gin.Context) {
//	var (
//		request = &models.MTBetOnNumReq{}
//		data    = make(map[string]interface{})
//	)
//
//	err := c.ShouldBindJSON(&request)
//	if err != nil {
//		controllers.Response(c, common.ParameterIllegal, "", data)
//		return
//	}
//
//	userId, err := helper.GetUserID(c)
//	if err != nil {
//		controllers.Response(c, common.UnauthorizedUserID, "", data)
//		return
//	}
//
//	reqData := &pbs.BetNumReq{
//		UserId: userId,
//		Bet:    int32(request.Bet),
//	}
//	reqMarshal, _ := proto.Marshal(reqData)
//
//	msgReq := pbs.NetMessage{
//		ReqHead: &pbs.ReqHead{
//			Uid:      0,
//			Token:    "",
//			Platform: "",
//		},
//		AckHead:   &pbs.AckHead{},
//		ServiceId: "",
//		MsgId:     int32(pbs.Meb_pnBetNumReq),
//		Content:   reqMarshal,
//	}
//	response, err := grpcclient.GetMtClient().CallMtMethod(&msgReq)
//	if response != nil && response.AckHead.Code != pbs.Code_OK {
//		controllers.Response(c, uint32(response.AckHead.Code), "", data)
//		return
//	}
//	if err != nil || response == nil {
//		global.GVA_LOG.Error("could not call method:", zap.Error(err))
//		controllers.Response(c, common.RpcCallRespErr, "", data)
//		return
//	}
//
//	data = gin.H{}
//	global.GVA_LOG.Infof("MtBetOnNum data:%v", data)
//	controllers.Response(c, common.WebOK, "", data)
//}
//
//// MtUserPeriodLayerList 用户（维度）每期每层历史
//// @Summary      用户（维度）每期每层历史
//// @Tags         魔塔
//// @Description  用户（维度）每期每层历史
//// @Accept       json
//// @Produce      json
//// @Param        user  body      models.MtUserPeriodLayerListReq                true  "用户（维度）每期每层历史"
//// @Success  1   {object}        common.JSONResult{data=models.MtUserPeriodLayerListResp} "用户（维度）每期每层历史"
//// @Failure      400   {object}  common.JSONResult                "错误提示"
//// @Router       /userPeriodLayerList [post]
//func MtUserPeriodLayerList(c *gin.Context) {
//	var (
//		request = &models.MtUserPeriodLayerListReq{}
//		data    = make(map[string]interface{})
//	)
//
//	err := c.ShouldBindJSON(&request)
//	if err != nil {
//		controllers.Response(c, common.ParameterIllegal, "", data)
//		return
//	}
//
//	userId, err := helper.GetUserID(c)
//	if err != nil {
//		controllers.Response(c, common.UnauthorizedUserID, "", data)
//		return
//	}
//	if request.PeriodId <= 0 {
//		controllers.Response(c, common.ParameterIllegal, "", data)
//		return
//	}
//
//	reqData := &pbs.UserPeriodLayerListReq{
//		UserId:   userId,
//		PeriodId: int32(request.PeriodId),
//	}
//	reqMarshal, _ := proto.Marshal(reqData)
//
//	msgReq := pbs.NetMessage{
//		ReqHead: &pbs.ReqHead{
//			Uid:      0,
//			Token:    "",
//			Platform: "",
//		},
//		AckHead:   &pbs.AckHead{},
//		ServiceId: "",
//		MsgId:     int32(pbs.Mmb_pnUserPeriodLayerListReq),
//		Content:   reqMarshal,
//	}
//	response, err := grpcclient.GetMtClient().CallMtMethod(&msgReq)
//	if response != nil && response.AckHead.Code != pbs.Code_OK {
//		controllers.Response(c, uint32(response.AckHead.Code), "", data)
//		return
//	}
//	if err != nil || response == nil {
//		global.GVA_LOG.Error("could not call method:", zap.Error(err))
//		controllers.Response(c, common.RpcCallRespErr, "", data)
//		return
//	}
//
//	userPeriodLayerListAck := pbs.UserPeriodLayerListAck{
//		UserPeriodLayerList: make([]*pbs.UserPeriodLayer, 0),
//	}
//	respData := response.Content
//	err = proto.Unmarshal(respData, &userPeriodLayerListAck)
//	if err != nil {
//		global.GVA_LOG.Error("Unmarshal mTStatusAck :", zap.Error(err))
//		controllers.Response(c, common.RpcCallRespDataErr, "", data)
//		return
//	}
//
//	global.GVA_LOG.Infof("mTStatusAck: %v", &userPeriodLayerListAck)
//
//	data = gin.H{
//		"user_period_layer_list": userPeriodLayerListAck.UserPeriodLayerList,
//	}
//	global.GVA_LOG.Infof("MtBetOnNum data:%v", data)
//	controllers.Response(c, common.WebOK, "", data)
//}
//
//// MtUserRevenueRank 用户每期的收益排行
//// @Summary      用户每期的收益排行
//// @Tags         魔塔
//// @Description  用户每期的收益排行
//// @Accept       json
//// @Produce      json
//// @Param        user  body      models.UserPeriodRevenueRankReq                true  "用户每期的收益排行"
//// @Success  1   {object}        common.JSONResult{data=models.UserPeriodRevenueRankResp} "用户每期的收益排行"
//// @Failure      400   {object}  common.JSONResult                "错误提示"
//// @Router       /userRevenueRank [post]
//func MtUserRevenueRank(c *gin.Context) {
//	var (
//		request = &models.UserPeriodRevenueRankReq{}
//		data    = make(map[string]interface{})
//	)
//
//	err := c.ShouldBindJSON(&request)
//	if err != nil {
//		controllers.Response(c, common.ParameterIllegal, "", data)
//		return
//	}
//
//	userId, err := helper.GetUserID(c)
//	if err != nil {
//		controllers.Response(c, common.UnauthorizedUserID, "", data)
//		return
//	}
//	if request.PeriodId <= 0 {
//		controllers.Response(c, common.ParameterIllegal, "", data)
//		return
//	}
//
//	reqData := &pbs.UserRevenueRankReq{
//		UserId:   userId,
//		PeriodId: int32(request.PeriodId),
//	}
//	reqMarshal, _ := proto.Marshal(reqData)
//
//	msgReq := pbs.NetMessage{
//		ReqHead: &pbs.ReqHead{
//			Uid:      0,
//			Token:    "",
//			Platform: "",
//		},
//		AckHead:   &pbs.AckHead{},
//		ServiceId: "",
//		MsgId:     int32(pbs.Mmb_pnUserRevenueRankReq),
//		Content:   reqMarshal,
//	}
//	response, err := grpcclient.GetMtClient().CallMtMethod(&msgReq)
//	if response != nil && response.AckHead.Code != pbs.Code_OK {
//		controllers.Response(c, uint32(response.AckHead.Code), "", data)
//		return
//	}
//	if err != nil || response == nil {
//		global.GVA_LOG.Error(" MtUserRevenueRank could not call method:", zap.Error(err))
//		controllers.Response(c, common.RpcCallRespErr, "", data)
//		return
//	}
//
//	userRevenueRankAck := pbs.UserRevenueRankAck{
//		UserRevenueRank: make([]*pbs.UserRevenue, 0),
//	}
//
//	respData := response.Content
//	err = proto.Unmarshal(respData, &userRevenueRankAck)
//	if err != nil {
//		global.GVA_LOG.Error("MtUserRevenueRank Unmarshal mTStatusAck :", zap.Error(err))
//		controllers.Response(c, common.RpcCallRespDataErr, "", data)
//		return
//	}
//
//	global.GVA_LOG.Infof("MtUserRevenueRank: %v", &userRevenueRankAck)
//
//	data = gin.H{
//		"user_period_revenue_ranks": userRevenueRankAck.UserRevenueRank,
//	}
//	global.GVA_LOG.Infof("MtBetOnNum data:%v", data)
//	controllers.Response(c, common.WebOK, "", data)
//}
//
//// MtOneTouchAddBetNum  一键召唤：充值
//// @Summary      		一键召唤充值
//// @Tags        		 魔塔
//// @Description  		一键召唤充值
//// @Accept       json
//// @Produce      json
//// @Param        user  body      models.MTOneTouchAddBetNumReq                true  "一键召唤充值"
//// @Success  1   {object}        common.JSONResult{data=models.MTOneTouchAddBetNumResp} "一键召唤充值"
//// @Failure      400   {object}  common.JSONResult                "错误提示"
//// @Router       /oneTouchBetNum [post]
//func MtOneTouchAddBetNum(c *gin.Context) {
//	var (
//		request = &models.MTOneTouchAddBetNumReq{}
//		data    = make(map[string]interface{})
//	)
//
//	err := c.ShouldBindJSON(&request)
//	if err != nil {
//		controllers.Response(c, common.ParameterIllegal, "", data)
//		return
//	}
//
//	userId, err := helper.GetUserID(c)
//	if err != nil {
//		controllers.Response(c, common.UnauthorizedUserID, "", data)
//		return
//	}
//
//	//本地锁
//	cacheKey := userId + "AutoNum"
//	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
//	if err != nil {
//		global.GVA_LOG.Infof("自动用户 || 一键召唤 QueueDataKeyMap TryAdd%v", cacheKey)
//		controllers.Response(c, common.DuplicateRequests, "", data)
//		return
//	}
//	defer global.QueueDataKeyMap.Del(cacheKey)
//
//	reqData := &pbs.BetOneTouchAddBetNumReq{
//		UserId: userId,
//		Bet:    int32(request.Bet),
//	}
//	reqMarshal, _ := proto.Marshal(reqData)
//
//	msgReq := pbs.NetMessage{
//		ReqHead: &pbs.ReqHead{
//			Uid:      0,
//			Token:    "",
//			Platform: "",
//		},
//		AckHead:   &pbs.AckHead{},
//		ServiceId: "",
//		//MsgId:     int32(pbs.Mmb_mtOneTouchAddBetReq),
//		Content:   reqMarshal,
//	}
//	response, err := grpcclient.GetMtClient().CallMtMethod(&msgReq)
//	if response != nil && response.AckHead.Code != pbs.Code_OK {
//		controllers.Response(c, uint32(response.AckHead.Code), "", data)
//		return
//	}
//	if err != nil || response == nil {
//		global.GVA_LOG.Error("MtOneTouchBetNum could not call method:", zap.Error(err))
//		controllers.Response(c, common.RpcCallRespErr, "", data)
//		return
//	}
//
//	data = gin.H{}
//	global.GVA_LOG.Infof("MtOneTouchBetNum MtBetOnNum data:%v", data)
//	controllers.Response(c, common.WebOK, "", data)
//}
//
//// MtJoinGameSelectCamp 选择阵营-新进入用户
//// @Summary      选择阵营-新进入用户
//// @Tags         魔塔
//// @Description  选择阵营-新进入用户
//// @Accept       json
//// @Produce      json
//// @Param        user  body      models.MTJoinGameSelectCampReq                true  "选择阵营-新进入用户"
//// @Success  1   {object}        common.JSONResult{data=models.MTJoinGameSelectCampResp} "选择阵营-新进入用户"
//// @Failure      400   {object}  common.JSONResult                "错误提示"
//// @Router       /JoinGameSelectCamp [post]
//func MtJoinGameSelectCamp(c *gin.Context) {
//	var (
//		request = &models.MTJoinGameSelectCampReq{}
//		data    = make(map[string]interface{})
//	)
//
//	err := c.ShouldBindJSON(&request)
//	if err != nil {
//		controllers.Response(c, common.ParameterIllegal, "", data)
//		return
//	}
//
//	userId, err := helper.GetUserID(c)
//	if err != nil {
//		controllers.Response(c, common.UnauthorizedUserID, "", data)
//		return
//	}
//
//	if request.Camp != 1 && request.Camp != 2 {
//		controllers.Response(c, common.ParameterIllegal, "", data)
//		return
//	}
//
//	reqData := &pbs.JoinGameSelectCampReq{
//		UserId: userId,
//		Camp:   int32(request.Camp),
//	}
//	reqMarshal, _ := proto.Marshal(reqData)
//
//	msgReq := pbs.NetMessage{
//		ReqHead: &pbs.ReqHead{
//			Uid:      0,
//			Token:    "",
//			Platform: "",
//		},
//		AckHead:   &pbs.AckHead{},
//		ServiceId: "",
//		MsgId:     int32(pbs.Mmb_pnJoinGameSelectCampReq),
//		Content:   reqMarshal,
//	}
//	response, err := grpcclient.GetMtClient().CallMtMethod(&msgReq)
//	if response != nil && response.AckHead.Code != pbs.Code_OK {
//		controllers.Response(c, uint32(response.AckHead.Code), "", data)
//		return
//	}
//	if err != nil || response == nil {
//		global.GVA_LOG.Error("MtJoinGameSelectCamp could not call method:", zap.Error(err))
//		controllers.Response(c, common.RpcCallRespErr, "", data)
//		return
//	}
//
//	data = gin.H{}
//	global.GVA_LOG.Infof("MtJoinGameSelectCamp MtBetOnNum data:%v", data)
//	controllers.Response(c, common.WebOK, "", data)
//}
