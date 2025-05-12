package websocket

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models"

	"sort"
	"time"
)

func (trs *RoomSpace) RobotAction() {
	if !trs.IsAllLoadComps {
		return
	}
	//当前在那个阶段
	//随牌|出牌|点赞
	//0=游戏未开始
	//1=游戏开始但是没有加载完成
	//2=用户随牌阶段
	//3=用户出牌阶段
	//4=用户点赞阶段
	//5=点赞界面 等待结算或者进入下一轮
	gameStatus, timeDown := trs.CurrGameTurnStateAndDownTime()
	global.GVA_LOG.Infof("机器人行为 RobotAction 当前的游戏状态:%v ，倒计时 gameStatus：%v", gameStatus, timeDown)
	if gameStatus == 0 || gameStatus == 1 {
		return
	}

	robotClassMap := make(map[int]int)
	firstUserMap := make(map[string]string)
	for _, userInfo := range trs.ComRoomSpace.UserInfos {
		_, ok := robotClassMap[userInfo.UserProperty.RobotClass]
		if ok {
			robotClassMap[userInfo.UserProperty.RobotClass]++
		} else {
			robotClassMap[userInfo.UserProperty.RobotClass] = 1
			firstUserMap[userInfo.UserID] = userInfo.UserID
		}
	}

	for _, userInfo := range trs.ComRoomSpace.UserInfos {
		if !userInfo.UserIsRobot() {
			continue
		}

		if robotClassMap[userInfo.UserProperty.RobotClass] > 1 {
			_, isFirst := firstUserMap[userInfo.UserID]
			if !isFirst {
				time.Sleep(1 * time.Second)
			}
		}

		//如果在出牌阶段 根基机器人的类型出牌
		if gameStatus == int(CliOutCard) {
			//出牌阶段
			if timeDown > 15 && timeDown <= 20 {
				if userInfo.UserProperty.RobotClass == 1 {
					//出牌
					global.GVA_LOG.Infof("机器人出牌 行为1 %v", userInfo.GetString())
					trs.Class1OutCardAction(userInfo)
				}
			}

			if timeDown > 20 && timeDown <= 25 {

				if userInfo.UserProperty.RobotClass == 2 {
					//出牌
					global.GVA_LOG.Infof("机器人出牌 行为2  %v", userInfo.GetString())
					trs.Class2OutCardAction(userInfo)
				}
			}

			if timeDown > 25 && timeDown <= 30 {
				if userInfo.UserProperty.RobotClass == 3 {
					//出牌
					global.GVA_LOG.Infof("机器人出牌 行为3  %v", userInfo.GetString())
					trs.Class3OutCardAction(userInfo)
				}
			}

		}

		if gameStatus == int(CliLikePage) {
			if timeDown > 1 && timeDown <= 15 {
				//if timeDown > 10 && timeDown <= 15 {
				if userInfo.UserProperty.RobotClass == 3 {
					//点赞
					//trs.Class3LikedAction(userInfo, notRobotUserIdArr)
					global.GVA_LOG.Infof("机器人点赞 行为3  %v", userInfo.GetString())
					trs.ClassStrategyLikedAction(userInfo, 3)
				}
			}

			//点赞阶段
			if timeDown > 0 && timeDown <= 5 {
				if userInfo.UserProperty.RobotClass == 1 {
					//点赞  strategy : 1:根据等级排序 2 随机
					global.GVA_LOG.Infof("机器人点赞 行为1  %v", userInfo.GetString())
					trs.ClassStrategyLikedAction(userInfo, 1)
				}
			}

			if timeDown > 5 && timeDown <= 10 {
				if userInfo.UserProperty.RobotClass == 2 {
					//点赞 strategy : 1:根据等级排序 2 随机
					global.GVA_LOG.Infof("机器人点赞 行为2  %v", userInfo.GetString())
					trs.ClassStrategyLikedAction(userInfo, 2)
				}
				//if userInfo.UserProperty.RobotClass == 3 {
				//	//点赞
				//}
			}
		}
	}
}

func (trs *RoomSpace) Class1OutCardAction(userInfo *models.UserInfo) {
	//帮助用户出牌，在手牌中随机抽牌打出，不设置明显偏好。

	//是否出牌
	cards := trs.ComRoomSpace.GetUserOutEdCards(userInfo.UserID)
	if len(cards) > 0 {
		//用户已经出过牌
		return
	}
	//获取用户当前的牌
	currCards, err := trs.ComRoomSpace.GetCurrCard(userInfo.UserID)
	if err != nil {
		global.GVA_LOG.Error("Class1OutCardAction 机器人1 GetCurrCard ", zap.Error(err))
	}
	if len(currCards) <= 0 {
		global.GVA_LOG.Error("Class1OutCardAction 机器人1 GetCurrCard == 0")
		return
	}
	var reqCards []*models.Card
	reqCards = append(reqCards, currCards[helper.RandInt(len(currCards))])
	global.GVA_LOG.Infof("Class1OutCardAction 机器人1 出牌 reqCards %v ,currCards  %v ,userInfo.UserID  %v", reqCards, currCards, userInfo.UserID)
	trs.OutCart(reqCards, currCards, userInfo.UserID) //机器人1 出牌
}

func (trs *RoomSpace) Class2OutCardAction(userInfo *models.UserInfo) {
	//在3-6秒内出牌，优先出最高等级的牌，同样等级的牌随机一张出牌

	//是否出牌
	cards := trs.ComRoomSpace.GetUserOutEdCards(userInfo.UserID)
	if len(cards) > 0 {
		//用户已经出过牌
		return
	}
	//获取用户当前的牌
	currCards, err := trs.ComRoomSpace.GetCurrCard(userInfo.UserID)
	if err != nil {
		global.GVA_LOG.Error("Class2OutCardAction 机器人2 GetCurrCard ", zap.Error(err))
	}
	if len(currCards) <= 0 {
		global.GVA_LOG.Error("Class2OutCardAction 机器人2 GetCurrCard == 0")
		return
	}

	// 使用 sort.Slice 实现倒序排序
	sort.Slice(currCards, func(i, j int) bool {
		// ">" 表示降序
		return currCards[i].Level > currCards[i].Level
	})

	var reqCards []*models.Card
	reqCards = append(reqCards, currCards[0])
	global.GVA_LOG.Infof("Class2OutCardAction 机器人2 出牌 reqCards %v ,currCards  %v ,userInfo.LikeUserId  %v", reqCards, currCards, userInfo.UserID)
	trs.OutCart(reqCards, currCards, userInfo.UserID) //机器人2 出牌
}

func (trs *RoomSpace) Class3OutCardAction(userInfo *models.UserInfo) {
	//在6-10秒内出牌，出牌优先级随机。
	//是否出牌
	cards := trs.ComRoomSpace.GetUserOutEdCards(userInfo.UserID)
	if len(cards) > 0 {
		//用户已经出过牌
		return
	}
	//获取用户当前的牌
	currCards, err := trs.ComRoomSpace.GetCurrCard(userInfo.UserID)
	if err != nil {
		global.GVA_LOG.Error("Class3OutCardAction 机器人3 GetCurrCard ", zap.Error(err))
	}
	if len(currCards) <= 0 {
		global.GVA_LOG.Error("Class3OutCardAction 机器人3 GetCurrCard == 0")
		return
	}

	var reqCards []*models.Card
	reqCards = append(reqCards, currCards[helper.RandInt(len(currCards))])
	global.GVA_LOG.Infof("Class3OutCardAction 机器人3 出牌 reqCards %v ,currCards  %v ,userInfo.UserID  %v", reqCards, currCards, userInfo.UserID)
	trs.OutCart(reqCards, currCards, userInfo.UserID) //机器人3 出牌
}

// ClassStrategyLikedAction  strategy: (1:根据等级排序 2:随机 3:跟风)
// 1:优先自动点赞场上品质最高的卡牌。如果场上出现相同品质的卡牌，则随机投一张卡牌。
// 2:完全随机点赞。点赞时间控制在2-7s随机分布。
// 3:优先跟票制，根据第一个点赞的场上玩家的点赞卡牌，在1-3秒内对其点赞进行跟票。（如果场上玩家都没有投票，则随机一张卡牌进行投票。
func (trs *RoomSpace) ClassStrategyLikedAction(userInfo *models.UserInfo, strategy int) {
	userId := userInfo.UserID
	likeUserId := ""
	likeCard := models.LikeCard{}
	outCards := make([]*models.Card, 0)
	likeCards := make([]*models.Card, 0)
	outLikeCards := make([]*models.LikeCard, 0)

	//用户是否点赞
	likeUserInfo := trs.ComRoomSpace.GetLikeUserInfo(userInfo.UserID)
	if len(likeUserInfo) > 0 {
		//该用户已经给别人点过赞
		global.GVA_LOG.Infof("Class1LikedAction 该用户已经给别人点过赞 userID %v", userInfo.UserID)
		return
	}

	//每个人出一个牌 取第一个就行
	outCards = trs.ComRoomSpace.GetUserOutEdCardExcludeUser(userInfo.UserID)
	if len(outCards) <= 0 {
		return
	}

	if strategy == 1 {
		// 使用 sort.Slice 实现倒序排序
		sort.Slice(outCards, func(i, j int) bool {
			return outCards[i].Level > outCards[j].Level
		})
	}
	if strategy == 2 {
		//完全随机
		helper.SliceShuffle(outCards)
	}

	if strategy == 1 || strategy == 2 {
		//找到牌等级最高的一张 点赞
		//每轮每次只出一张牌
		isOutCard := false
		for _, outCard := range outCards {
			if isOutCard {
				//只取第0个
				break
			}
			likeCard = models.LikeCard{
				CardId:     outCard.CardId,
				LikeUserId: outCard.UserID,
				Level:      outCard.Level,
				AddRate:    outCard.AddRate,
			}
			isOutCard = true
			likeUserId = outCard.UserID
			likeCards = append(likeCards, outCard)
		}
		if !isOutCard {
			return
		}
		trs.DoLikeCard(userId, likeUserId, likeCard, likeCards)
	}

	if strategy == 3 {
		notSelfUserIdArr := make([]string, 0)
		for _, uInfo := range trs.ComRoomSpace.UserInfos {
			if uInfo.UserID == userId {
				continue
			}
			notSelfUserIdArr = append(notSelfUserIdArr, uInfo.UserID)
		}

		//机器人跟风对象
		followLikeUserId := ""
		for _, notRobotUserId := range notSelfUserIdArr {
			robotLikeUserInfo := trs.ComRoomSpace.GetLikeUserInfo(notRobotUserId)
			if len(robotLikeUserInfo) <= 0 {
				continue
			}
			likeSelf := false
			for kk, _ := range robotLikeUserInfo {
				//如果给自己点赞 跳过
				if robotLikeUserInfo[kk].LikeUserId == userId {
					likeSelf = true
					break
				}
			}
			if likeSelf {
				continue
			}

			//每次只能点赞一张牌
			followLikeUserId = robotLikeUserInfo[0].LikeUserId
			outLikeCards = robotLikeUserInfo
			break
		}
		if len(followLikeUserId) <= 0 {
			return
		}

		isLikeCard := false
		for _, outLikeCard := range outLikeCards {
			if isLikeCard == true {
				break
			}
			if outLikeCard.LikeUserId == userId {
				continue
			}
			likeCard = models.LikeCard{
				CardId:     outLikeCard.CardId,
				LikeUserId: outLikeCard.LikeUserId,
				Level:      outLikeCard.Level,
				AddRate:    outLikeCard.AddRate,
			}
			isLikeCard = true
			likeUserId = outLikeCard.LikeUserId
			likeCards = trs.ComRoomSpace.GetUserOutEdCards(outLikeCard.LikeUserId)
		}
		if !isLikeCard {
			return
		}

		time.Sleep(1 * time.Second)
		//global.GVA_LOG.Infof("outLikeCards%v likeCards:%v", outLikeCards, likeCards)
		trs.DoLikeCard(userId, likeUserId, likeCard, likeCards)
	}
}

//func (trs *RoomSpace) Class3LikedAction(userInfo *models.UserInfo, notRobotUserIdArr []string) {
//	//优先跟票制，根据第一个点赞的场上玩家的点赞卡牌，在1-3秒内对其点赞进行跟票。（如果场上玩家都没有投票，则随机一张卡牌进行投票。）
//	if len(notRobotUserIdArr) <= 0 {
//		return
//	}
//
//	//判断有没有点赞 如果有人点赞就跟
//	userId := userInfo.UserID
//	likeUserId := ""
//	likeCard := models.LikeCard{}
//	likeCards := make([]*models.Card, 0)
//	outLikeCards := make([]*models.LikeCard, 0)
//
//	//用户是否点赞
//	likeUserInfo := trs.ComRoomSpace.GetLikeUserInfo(userInfo.UserID)
//	if len(likeUserInfo) > 0 {
//		//该用户已经给别人点过赞
//		global.GVA_LOG.Infof("Class1LikedAction 该用户已经给别人点过赞 userID %v", userInfo.UserID)
//		return
//	}
//
//	//机器人跟风对象
//	followLikeUserId := ""
//	for _, notRobotUserId := range notRobotUserIdArr {
//		robotLikeUserInfo := trs.ComRoomSpace.GetLikeUserInfo(notRobotUserId)
//		if len(robotLikeUserInfo) > 0 {
//			followLikeUserId = notRobotUserId
//			outLikeCards = robotLikeUserInfo
//			break
//		}
//	}
//	if len(followLikeUserId) <= 0 {
//		return
//	}
//
//	isLikeCard := false
//	for _, outLikeCard := range outLikeCards {
//		if isLikeCard == true {
//			break
//		}
//		if outLikeCard.LikeUserId == userId {
//			continue
//		}
//		likeCard = models.LikeCard{
//			CardId:     outLikeCard.CardId,
//			LikeUserId: outLikeCard.LikeUserId,
//			Level:      outLikeCard.Level,
//			AddRate:    outLikeCard.AddRate,
//		}
//		isLikeCard = true
//		likeUserId = outLikeCard.LikeUserId
//		likeCards = trs.ComRoomSpace.GetUserOutEdCards(outLikeCard.LikeUserId)
//	}
//	if !isLikeCard {
//		return
//	}
//
//	time.Sleep(1 * time.Second)
//
//	//global.GVA_LOG.Infof("outLikeCards%v likeCards:%v", outLikeCards, likeCards)
//	trs.DoLikeCard(userId, likeUserId, likeCard, likeCards)
//}

func (trs *RoomSpace) Class1ReMakeAction(userInfo *models.UserInfo) {
	//不执行任何操作。
}

func (trs *RoomSpace) Class2ReMakeAction(userInfo *models.UserInfo) {
	//随机进行1-2次重随。

}

func (trs *RoomSpace) Class3ReMakeAction(userInfo *models.UserInfo) {
	//随机执行2-3次重随
}

func (trs *RoomSpace) Class1LikedAction(userInfo *models.UserInfo) {
	//优先自动点赞场上品质最高的卡牌。如果场上出现相同品质的卡牌，则随机投一张卡牌。

	userId := userInfo.UserID
	likeUserId := ""
	likeCard := models.LikeCard{}
	outCards := make([]*models.Card, 0)
	likeCards := make([]*models.Card, 0)

	//那个用户没点赞
	likeUserInfo := trs.ComRoomSpace.GetLikeUserInfo(userInfo.UserID)
	if len(likeUserInfo) > 0 {
		//该用户已经给别人点过赞
		global.GVA_LOG.Infof("Class1LikedAction 该用户已经给别人点过赞 userID %v", userInfo.UserID)
		return
	}

	//每个人出一个牌 取第一个就行
	outCards = trs.ComRoomSpace.GetUserOutEdCardExcludeUser(userInfo.UserID)
	if len(outCards) <= 0 {
		return
	}

	// 使用 sort.Slice 实现倒序排序
	sort.Slice(outCards, func(i, j int) bool {
		return outCards[i].Level > outCards[j].Level
	})

	//找到牌等级最高的一张 点赞
	//每轮每次只出一张牌
	isOutCard := false
	for _, outCard := range outCards {
		if isOutCard == true {
			break
		}
		//if outCard.UserID == userId {
		//	continue
		//}
		likeCard = models.LikeCard{
			CardId:     outCard.CardId,
			LikeUserId: outCard.UserID,
			Level:      outCard.Level,
			AddRate:    outCard.AddRate,
		}
		isOutCard = true
		likeUserId = outCard.UserID
		likeCards = append(likeCards, outCard)
	}

	if !isOutCard {
		return
	}
	trs.DoLikeCard(userId, likeUserId, likeCard, likeCards)
}

//func (trs *RoomSpace) Class2LikedAction(userInfo *models.UserInfo) {
//	//完全随机点赞。点赞时间控制在2-7s随机分布。
//	userId := userInfo.LikeUserId
//	likeUserId := ""
//	likeCard := models.LikeCard{}
//	outCards := make([]*models.Card, 0)
//	likeCards := make([]*models.Card, 0)
//
//	//那个用户没点赞
//	likeUserInfo := trs.ComRoomSpace.GetLikeUserInfo(userInfo.LikeUserId)
//	if len(likeUserInfo) > 0 {
//		//该用户已经给别人点过赞
//		global.GVA_LOG.Infof("Class1LikedAction 该用户已经给别人点过赞 userID %v", userInfo.LikeUserId)
//		return
//	}
//
//	//每个人出一个牌 取第一个就行
//	outCards = trs.ComRoomSpace.GetUserOutEdCardExcludeUser(userInfo.LikeUserId)
//	if len(outCards) <= 0 {
//		return
//	}
//
//	//完全随机
//	helper.SliceShuffle(outCards)
//
//	//找到牌等级最高的一张 点赞
//	//每轮每次只出一张牌
//	isOutCard := false
//	for k, outCard := range outCards {
//		if k == 1 {
//			//只取第0个
//			break
//		}
//		likeCard = models.LikeCard{
//			CardId:  outCard.CardId,
//			LikeUserId:  outCard.LikeUserId,
//			Level:   outCard.Level,
//			AddRate: outCard.AddRate,
//		}
//		isOutCard = true
//		likeUserId = outCard.LikeUserId
//		likeCards = append(likeCards, outCard)
//	}
//
//	if !isOutCard {
//		return
//	}
//	trs.DoLikeCard(userId, likeUserId, likeCard, likeCards)
//}

//
//func (trs *RoomSpace) ByRobotClassSetAction(uInfo *models.UserInfo) {
//	if !uInfo.UserIsRobot() {
//		return
//	}
//
//	//用户维度
//
//	if uInfo.UserProperty.RobotClass == 1 {
//
//	}
//
//	if uInfo.UserProperty.RobotClass == 2 {
//
//	}
//
//	if uInfo.UserProperty.RobotClass == 3 {
//
//	}
//
//	////设置随牌倒计时
//	//uInfo.SetReMakeCardDown(models.GetReMakeCardDownTimeInt(CommTimeOut))
//	////出牌倒计时
//	//uInfo.SetOutCardCountDown(models.GetOutCardCountDownTimeInt(CommTimeOutDouble))
//	////设置重随牌状态
//	//uInfo.SetGameStatus(int(RemakeCardIng))
//}
