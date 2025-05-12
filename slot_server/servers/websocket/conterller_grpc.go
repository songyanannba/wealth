package websocket

import (
	"context"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"slot_server/lib/common"
	"slot_server/lib/config"
	"slot_server/lib/global"
	"slot_server/lib/models/table"
	"slot_server/lib/src/component"
	"slot_server/lib/src/dao"
	"slot_server/lib/src/logic"
	"slot_server/protoc/pbs"
)

// friendUserList
func FriendUserList(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
	global.GVA_LOG.Infof("好友列表 FriendUserList：%v", req)
	var (
		request       = pbs.FriendListReq{}
		comResp       = component.NewNetMessage(int32(pbs.Meb_friendUserListResp))
		listData      = []*pbs.UserFriend{}
		friendListAck = &pbs.FriendListAck{}
	)

	if req.MsgId != int32(pbs.Meb_friendUserList) {
		global.GVA_LOG.Error("FriendUserList 协议号不正确", zap.Any("FriendUserList", req))
		comResp.AckHead.Code = pbs.Code_ProtocNumberError
		return comResp, nil
	}
	err := proto.Unmarshal(req.Content, &request)
	if err != nil {
		comResp.AckHead.Code = pbs.Code_DataCompileError
		return comResp, nil
	}
	global.GVA_LOG.Infof("获取用户矿工列表 UserRatList：%v", &request)

	userFriendLists, isHaveNextPage := logic.UserFriendAndIsHaveNextPage(request.UserId, int(request.LastId))

	for _, userFriendList := range userFriendLists {
		listData = append(listData, &pbs.UserFriend{
			FriendUserId: userFriendList.FriendUserId,
			Nickname:     userFriendList.Nickname,
			FriendId:     int32(userFriendList.FriendId),
		})
	}

	friendListAck = &pbs.FriendListAck{
		IsHaveNextPage: isHaveNextPage,
		UserFriend:     listData,
	}

	//global.GVA_LOG.Infof("FriendUserList 获取用户矿工列表：%v", friendListAck)
	//返回数据
	friendListAckMarshal, _ := proto.Marshal(friendListAck)
	comResp.Content = friendListAckMarshal
	return comResp, nil
}

// AuditUserList
func AuditUserList(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
	global.GVA_LOG.Infof("好友审核列表 AuditUserList：%v", req)
	var (
		request       = pbs.AuditUserListReq{}
		comResp       = component.NewNetMessage(int32(pbs.Meb_auditUserListResp))
		listData      = []*pbs.AuditUser{}
		friendListAck = &pbs.AuditUserAck{}
	)

	if req.MsgId != int32(pbs.Meb_auditUserList) {
		global.GVA_LOG.Error("AuditUserList 协议号不正确", zap.Any("AuditUserList", req))
		comResp.AckHead.Code = pbs.Code_ProtocNumberError
		return comResp, nil
	}
	err := proto.Unmarshal(req.Content, &request)
	if err != nil {
		comResp.AckHead.Code = pbs.Code_DataCompileError
		return comResp, nil
	}
	global.GVA_LOG.Infof(" AuditUserList：%v", &request)

	authUserLists, isHaveNextPage := logic.AuditUserAndIsHaveNextPage(request.UserId, int(request.LastId))

	for _, userList := range authUserLists {
		listData = append(listData, &pbs.AuditUser{
			ApplicationUser: userList.ApplicationUser,
			Nickname:        userList.Nickname,
			AuditId:         int32(userList.AuditId),
		})
	}

	friendListAck = &pbs.AuditUserAck{
		IsHaveNextPage: isHaveNextPage,
		AuditUser:      listData,
	}

	//global.GVA_LOG.Infof("AuditUserList ：%v", friendListAck)
	//返回数据
	friendListAckMarshal, _ := proto.Marshal(friendListAck)
	comResp.Content = friendListAckMarshal
	return comResp, nil
}

func HandbookListController(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
	var (
		request         = pbs.HandbookListReq{}
		comResp         = component.NewNetMessage(int32(pbs.Meb_handbookListResp))
		listData        = make([]*pbs.HandListCard, 0)
		handbookListAck = &pbs.HandbookListAck{}
	)

	if req.MsgId != int32(pbs.Meb_handbookList) {
		global.GVA_LOG.Error("HandbookListController 协议号不正确", zap.Any("HandbookListController", req))
		comResp.AckHead.Code = pbs.Code_ProtocNumberError
		return comResp, nil
	}
	err := proto.Unmarshal(req.Content, &request)
	if err != nil {
		comResp.AckHead.Code = pbs.Code_DataCompileError
		return comResp, nil
	}
	global.GVA_LOG.Infof("HandbookListController：%v", &request)

	cardLists, isHaveNextPage, allCount := logic.CardConfListAndIsHaveNextPage(request.UserId, int(request.LastId), int(request.Level))

	for _, cardList := range cardLists {
		listData = append(listData, &pbs.HandListCard{
			CardId: int32(cardList.CardId),
			Name:   cardList.Name,
			Suffix: cardList.Suffix,
			IsOwn:  cardList.IsOwn,
			Level:  int32(cardList.Level),
		})
	}

	handbookListAck = &pbs.HandbookListAck{
		IsHaveNextPage: isHaveNextPage,
		AllCartCount:   int32(allCount),
		HandListCard:   listData,
	}

	global.GVA_LOG.Infof("handbookListAck：%v", handbookListAck)
	//返回数据
	friendListAckMarshal, _ := proto.Marshal(handbookListAck)
	comResp.Content = friendListAckMarshal
	return comResp, nil
}

func CardVersionListController(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
	var (
		request  = pbs.CardVersionListReq{}
		comResp  = component.NewNetMessage(int32(pbs.Meb_cardVersionListResp))
		listData = make([]*pbs.CardVersionList, 0)
	)

	if req.MsgId != int32(pbs.Meb_cardVersionList) {
		global.GVA_LOG.Error("CardVersionListController 协议号不正确", zap.Any("CardVersionListController", req))
		comResp.AckHead.Code = pbs.Code_ProtocNumberError
		return comResp, nil
	}
	err := proto.Unmarshal(req.Content, &request)
	if err != nil {
		comResp.AckHead.Code = pbs.Code_DataCompileError
		return comResp, nil
	}
	global.GVA_LOG.Infof(" CardVersionListController：%v", &request)

	cardLists := dao.GetCardConfVersion()

	for _, cardList := range cardLists {
		listData = append(listData, &pbs.CardVersionList{Version: int32(cardList.Version)})
	}

	ListAck := &pbs.CardVersionListAck{CardVersionList: listData}

	global.GVA_LOG.Infof("CardVersionListController UserID:{%v} 给客户端发消息:{%v}", request.UserId, ListAck)
	//返回数据
	listAckMarshal, _ := proto.Marshal(ListAck)
	comResp.Content = listAckMarshal
	return comResp, nil
}

func UnpackCardController(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
	var (
		request     = pbs.UnpackCardReq{}
		comResp     = component.NewNetMessage(int32(pbs.Meb_unpackCardResp))
		listDataArr = make([]*pbs.HandListCardArr, 0)
		listAck     = &pbs.UnpackCardAck{}
	)

	if req.MsgId != int32(pbs.Meb_unpackCard) {
		global.GVA_LOG.Error("UnpackCardController 协议号不正确", zap.Any("UnpackCardController", req))
		comResp.AckHead.Code = pbs.Code_ProtocNumberError
		return comResp, nil
	}
	err := proto.Unmarshal(req.Content, &request)
	if err != nil {
		comResp.AckHead.Code = pbs.Code_DataCompileError
		return comResp, nil
	}
	global.GVA_LOG.Infof("UnpackCardController：%v", &request)

	var isContinue bool
	//是否首次5连抽
	record, _ := table.UserIsHaveHandbook(request.UserId)
	if record.ID > 0 {
		isContinue = true
	}
	if isContinue && request.Num == 5 {
		comResp.AckHead.Code = pbs.Code_NotContinue5
		return comResp, nil
	}

	//判断金币是否够
	experienceInfo := dao.GetUserCoinExperience(request.UserId)
	//金币消耗配置
	coinConsumeConfig := dao.GetCoinConsumeConfigByType(int(table.CoinConsumeUnpackCard))
	if isContinue {
		if request.Num == config.UnpackCardNum1 && experienceInfo.CoinNum < coinConsumeConfig.CoinNum {
			comResp.AckHead.Code = pbs.Code_KingCoinNotEnough
			return comResp, nil
		}
	}

	unpackCardLists := logic.UnpackCardVersionAndNum(request.UserId, int(request.Version), int(request.Num))
	for _, cardListItem := range unpackCardLists {
		listData := make([]*pbs.HandListCard, 0)
		for _, cardList := range cardListItem {
			listData = append(listData, &pbs.HandListCard{
				CardId: int32(cardList.CardId),
				Name:   cardList.Name,
				Suffix: cardList.Suffix,
				Level:  int32(cardList.Level),
			})
		}
		listDataArr = append(listDataArr, &pbs.HandListCardArr{Cards: listData})
	}
	listAck.ListCard = listDataArr

	if len(unpackCardLists) > 0 && isContinue {
		//扣金币
		err := dao.UpdateUserCoinNumOrExperience(request.UserId, -coinConsumeConfig.CoinNum, 0, 1)
		if err != nil {
			global.GVA_LOG.Error("UnpackCardController ,UpdateUserCoinNumOrExperience err", zap.Any("err", err))
		}
	}

	//返回数据
	listAckMarshal, _ := proto.Marshal(listAck)
	comResp.Content = listAckMarshal
	return comResp, nil
}

func AddFriendController(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
	var (
		request = pbs.AddFriendReq{}
		comResp = component.NewNetMessage(int32(pbs.Meb_addFriendResp))
		ack     = &pbs.AddFriendAck{}
	)

	if req.MsgId != int32(pbs.Meb_addFriend) {
		global.GVA_LOG.Error("AddFriendController 协议号不正确", zap.Any("AddFriendController", req))
		comResp.AckHead.Code = pbs.Code_ProtocNumberError
		return comResp, nil
	}
	err := proto.Unmarshal(req.Content, &request)
	if err != nil {
		comResp.AckHead.Code = pbs.Code_DataCompileError
		return comResp, nil
	}
	global.GVA_LOG.Infof("AddFriendController：%v", &request)

	//判断是否已经在审核中 和 是否是好友
	applicationCode := logic.UserFriendAuditByAuditAndApplication(request.AuditUser, request.ApplicationUser)
	if applicationCode != common.OK {
		comResp.AckHead.Code = pbs.Code(applicationCode)
		return comResp, nil
	}

	//放入到审核列表
	logic.CreateUserFriendAudit(request.AuditUser, request.ApplicationUser)

	ackMarshal, _ := proto.Marshal(ack)
	comResp.Content = ackMarshal

	return comResp, nil
}

func DelFriendController(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
	var (
		request = pbs.DelFriendReq{}
		comResp = component.NewNetMessage(int32(pbs.Meb_delFriendResp))
		ack     = &pbs.DelFriendAck{}
	)

	if req.MsgId != int32(pbs.Meb_delFriend) {
		global.GVA_LOG.Error("DelFriendController 协议号不正确", zap.Any("DelFriendController", req))
		comResp.AckHead.Code = pbs.Code_ProtocNumberError
		return comResp, nil
	}
	err := proto.Unmarshal(req.Content, &request)
	if err != nil {
		comResp.AckHead.Code = pbs.Code_DataCompileError
		return comResp, nil
	}
	global.GVA_LOG.Infof("DelFriendController：%v", &request)

	record, err := table.GetUserFriendById(int(request.FriendId))
	if err != nil {
		comResp.AckHead.Code = pbs.Code_DBErr
		return comResp, nil
	}

	if record.ID == 0 {
		//没有申请纪录
		comResp.AckHead.Code = pbs.Code(common.NotFriendRecord)
		return comResp, nil
	}

	err = table.DelUserFriendById(int(request.FriendId))
	if err != nil {
		global.GVA_LOG.Error("DelUserFriendById", zap.Any("DelUserFriendById", request.FriendId))
		comResp.AckHead.Code = pbs.Code_DBErr
		return comResp, nil
	}

	//双向删除
	err = table.DelUserFriendByUserIdAndFriendUserId(record.FriendUserId, record.UserId)
	if err != nil {
		global.GVA_LOG.Error("", zap.Any("DelUserFriendByUserIdAndFriendUserId", record.UserId))
	}

	ackMarshal, _ := proto.Marshal(ack)
	comResp.Content = ackMarshal
	return comResp, nil
}

func AuthFriendController(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
	var (
		request = pbs.AuthFriendReq{}
		comResp = component.NewNetMessage(int32(pbs.Meb_authFriendResp))
		ack     = &pbs.AuthFriendAck{}
	)

	if req.MsgId != int32(pbs.Meb_authFriend) {
		global.GVA_LOG.Error("AuthFriendController 协议号不正确", zap.Any("AuthFriendController", req))
		comResp.AckHead.Code = pbs.Code_ProtocNumberError
		return comResp, nil
	}
	err := proto.Unmarshal(req.Content, &request)
	if err != nil {
		comResp.AckHead.Code = pbs.Code_DataCompileError
		return comResp, nil
	}
	global.GVA_LOG.Infof("AuthFriendController：%v", &request)

	//是否在审核列表 	//是否审核过

	application, err := table.UserFriendAuditById(int(request.AuditId))
	if err != nil {
		global.GVA_LOG.Error("UserFriendAuditByAuditAndApplication:", zap.Error(err))
		comResp.AckHead.Code = pbs.Code_DBErr
		return comResp, nil
	}
	if application.ID == 0 {
		//没有申请纪录
		comResp.AckHead.Code = pbs.Code(common.NotAuthRecord)
		return comResp, nil
	}

	if application.ID == 1 {
		//已经是好友
		comResp.AckHead.Code = pbs.Code(common.AuthHavePass)
		return comResp, nil
	}

	record, err := table.GetUserFriendByUserIdAndFriendId(application.AuditUser, application.ApplicationUser)
	if err != nil {
		global.GVA_LOG.Error("CreateUserFriend:", zap.Error(err))
	}
	if record.ID > 0 {
		comResp.AckHead.Code = pbs.Code(common.HaveFriend)
		return comResp, nil
	}

	//加到朋友列表 双向添加
	//logic.CreateUserFriend(application.AuditUser, application.ApplicationUser)
	//logic.CreateUserFriend(application.ApplicationUser, application.AuditUser)
	friendRecord, _ := table.GetUserFriendByUserIdAndFriendId(application.AuditUser, application.ApplicationUser)
	if friendRecord.ID <= 0 {
		logic.CreateUserFriend(application.AuditUser, application.ApplicationUser)
	}
	friendRecord, _ = table.GetUserFriendByUserIdAndFriendId(application.ApplicationUser, application.AuditUser)
	if friendRecord.ID <= 0 {
		logic.CreateUserFriend(application.AuditUser, application.ApplicationUser)
	}

	application.IsAgree = 1
	err = table.SaveUserFriendAudit(application)
	if err != nil {
		global.GVA_LOG.Error("CreateUserFriend:", zap.Error(err))
	}

	ackMarshal, _ := proto.Marshal(ack)
	comResp.Content = ackMarshal
	return comResp, nil
}

func UserDetailController(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
	var (
		request = pbs.UserDetailReq{}
		comResp = component.NewNetMessage(int32(pbs.Meb_userDetailResp))
		ack     = &pbs.UserDetailAck{}
	)

	if req.MsgId != int32(pbs.Meb_userDetail) {
		global.GVA_LOG.Error("UserDetailController 协议号不正确", zap.Any("UserDetailController", req))
		comResp.AckHead.Code = pbs.Code_ProtocNumberError
		return comResp, nil
	}
	err := proto.Unmarshal(req.Content, &request)
	if err != nil {
		comResp.AckHead.Code = pbs.Code_DataCompileError
		return comResp, nil
	}
	global.GVA_LOG.Infof("UserDetailController：%v", &request)

	userInfo, err := table.GetGameUserByUid(request.UserId)
	if err != nil {
		global.GVA_LOG.Error("UserDetailController GetGameUserByUid", zap.Error(err))
		comResp.AckHead.Code = pbs.Code_DBErr
		return comResp, nil
	}

	ack.UserId = userInfo.UserId
	ack.Nickname = userInfo.Nickname

	ackMarshal, _ := proto.Marshal(ack)
	comResp.Content = ackMarshal
	return comResp, nil
}

func CoinExperienceController(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
	var (
		request = pbs.CoinExperienceReq{}
		comResp = component.NewNetMessage(int32(pbs.Meb_coinExperienceResp))
		ack     = &pbs.CoinExperienceAck{}
	)

	if req.MsgId != int32(pbs.Meb_coinExperience) {
		global.GVA_LOG.Error("CoinExperienceController 协议号不正确", zap.Any("CoinExperienceController", req))
		comResp.AckHead.Code = pbs.Code_ProtocNumberError
		return comResp, nil
	}
	err := proto.Unmarshal(req.Content, &request)
	if err != nil {
		comResp.AckHead.Code = pbs.Code_DataCompileError
		return comResp, nil
	}
	global.GVA_LOG.Infof("CoinExperienceController：%v", &request)

	experienceInfo := dao.GetUserCoinExperience(request.UserId)

	ack.Experience = float32(experienceInfo.Experience)
	ack.CoinNum = float32(experienceInfo.CoinNum)

	ackMarshal, _ := proto.Marshal(ack)
	comResp.Content = ackMarshal
	return comResp, nil
}
