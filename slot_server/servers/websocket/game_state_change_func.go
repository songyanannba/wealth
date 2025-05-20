package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"slot_server/lib/config"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models"
	"slot_server/lib/src/dao"
	"slot_server/protoc/pbs"
	"sort"
	"strconv"
	"time"
)

func WheelAnimalSortCalculateExec(trs *RoomSpace) {
	global.GVA_LOG.Infof("WheelAnimalSortCalculateExec 房间 {%v}", trs.RoomInfo.RoomNo)
	var (
		msgData = &pbs.AnimalSortMsg{WinBetZoneConfig: make([]*pbs.WinBetZoneConfig, 0)}
	)

	//实际外部排序
	RecursionGetAnimalConfig(trs)

	//当前最外圈 动物的排序情况
	currAnimalWheelSort := trs.ComRoomSpace.CurrAnimalWheelSort

	for _, animalWheelSort := range currAnimalWheelSort {
		winBetZoneConfig := &pbs.WinBetZoneConfig{
			WinSeat: int32(animalWheelSort.WinSeat),
		}
		for _, animalConf := range animalWheelSort.AnimalConfigs {
			winBetZoneConfig.AnimalConfig = append(winBetZoneConfig.AnimalConfig, &pbs.AnimalConfig{
				Seat:     int32(animalConf.Seat),
				AnimalId: int32(animalConf.AnimalId),
			})
		}

		//押注大小
		for _, bigOrSmallConfig := range animalWheelSort.BigOrSmallConfigs {
			winBetZoneConfig.BigSmallConfig = append(winBetZoneConfig.BigSmallConfig, &pbs.BigOrSmallConfig{
				Seat:       int32(bigOrSmallConfig.Seat),
				BigSmallId: int32(bigOrSmallConfig.BigOrSmall),
			})
		}

		//对应位置的颜色
		colorConfigSeat := trs.GetColorConfigsBySeat(animalWheelSort.WinAnimalConfig.Seat)

		//根据本局赢钱的位置的动物和颜色确定赔率
		betZoneConfig := GetBetZoneConfigByAnimalIdAndColorId(animalWheelSort.WinAnimalConfig.AnimalId, colorConfigSeat.ColorId)
		animalWheelSort.WinBetZoneConfig = betZoneConfig

		for _, bzz := range betZoneConfig {
			winBetZoneConfig.WinZoneConf = append(winBetZoneConfig.WinZoneConf, &pbs.WinZoneConf{
				BetZoneId: int32(bzz.Seat),
				BetRate:   float32(bzz.BetRate),
			})
		}
		msgData.WinBetZoneConfig = append(msgData.WinBetZoneConfig, winBetZoneConfig)
	}

	//获取房间人数
	global.GVA_LOG.Infof("  押注停止后 主动下发最外圈的动物排序，第一个排在最上面 位置0开始: %v", msgData)
	responseHeadByte, _ := proto.Marshal(msgData)
	netMessageResp := helper.NewNetMessage("", "", int32(pbs.ProtocNum_AnimalSortMsg), config.SlotServer)
	netMessageResp.Content = responseHeadByte
	NatsSendAimUserMsg(trs, netMessageResp, "")
}

func RecursionGetAnimalConfig(trs *RoomSpace) {
	//最外圈的动物排序
	animalWheelSort := make([]*AllAnimalWheelSort, 0)
	//第一次获取最外圈的排序。0位置为顶点位置
	firstAnimalConfigs := WheelAnimalSortCalculate(trs)

	//todo
	//现在是随机
	//这个要根据押注情况，分析一个可以赢钱的位置
	firstAnimalConfigsLen := len(firstAnimalConfigs)
	winSeat := helper.RandInt(firstAnimalConfigsLen)
	//根据位置获取赢钱的动物
	winAnimalConfig := trs.GetNewAnimalConfigsBySeat(winSeat, firstAnimalConfigs)

	//大小只出现一次
	bigOrSmallConfigs := WheelBigOrSmallCalculate(trs)
	winBigOrSmallConfig := trs.GetBigOrSmallConfigsBySeat(winSeat, bigOrSmallConfigs)

	newAnimalWheelSort1 := &AllAnimalWheelSort{
		WinSeat:             winSeat,
		AnimalConfigs:       firstAnimalConfigs,
		WinAnimalConfig:     winAnimalConfig,
		BigOrSmallConfigs:   bigOrSmallConfigs,
		WinBigOrSmallConfig: winBigOrSmallConfig,
	}
	animalWheelSort = append(animalWheelSort, newAnimalWheelSort1)

	//如果是 LUCKY
	if winAnimalConfig.AnimalId == 2 {
		//先吧 第一次的 LUCKY， 放进结果集
		counts := []int{2, 3, 4, 5, 6, 7}
		//转几次
		randCountInt := helper.RandInt(len(counts))
		//todo
		randCountInt = 2

		for i := 0; i < randCountInt; i++ {
			//新的动物排序
			newAnimalConfigs := WheelAnimalSortCalculate(trs)
			//现在是随机
			winSeat = helper.RandInt(len(newAnimalConfigs))
			//根据位置获取动物
			winAnimalConfig = trs.GetNewAnimalConfigsBySeat(winSeat, newAnimalConfigs)

			//如果是 幸运的动物 位置+1
			if winAnimalConfig.AnimalId == 2 {
				winSeat += 1
				winAnimalConfig = trs.GetNewAnimalConfigsBySeat(winSeat, newAnimalConfigs)
				newAnimalWheelSort := &AllAnimalWheelSort{
					WinSeat:         winSeat,
					WinAnimalConfig: winAnimalConfig,
					AnimalConfigs:   newAnimalConfigs,
				}
				animalWheelSort = append(animalWheelSort, newAnimalWheelSort)
			} else {
				newAnimalWheelSort := &AllAnimalWheelSort{
					WinSeat:         winSeat,
					WinAnimalConfig: winAnimalConfig,
					AnimalConfigs:   newAnimalConfigs,
				}
				animalWheelSort = append(animalWheelSort, newAnimalWheelSort)
			}
		}
	}

	//当前的排序
	trs.ComRoomSpace.CurrAnimalWheelSort = animalWheelSort
}

func WheelAnimalSortCalculate(trs *RoomSpace) []*AnimalConfig {
	//todo 优化
	//要根据当前的押注 计算可以盈利的区间 然后指定到合适的位置
	animalConfigsLen := len(trs.AnimalConfigs)
	topSeat := helper.RandInt(animalConfigsLen)

	newAnimalConfigs := make([]*AnimalConfig, 0)
	newAnimalConfigs = append(newAnimalConfigs, trs.AnimalConfigs[topSeat:]...)
	newAnimalConfigs = append(newAnimalConfigs, trs.AnimalConfigs[:topSeat]...)

	for k, _ := range newAnimalConfigs {
		newAnimalConfigs[k].Seat = k
	}

	//trs.ComRoomSpace.CurrAnimalWheelSort = append(trs.ComRoomSpace.CurrAnimalWheelSort, newAnimalConfigs)
	return newAnimalConfigs
}

func WheelBigOrSmallCalculate(trs *RoomSpace) []*BigOrSmallConfig {

	//要根据当前的押注 计算可以盈利的区间 然后指定到合适的位置
	bigOrSmallConfigLen := len(trs.BigOrSmallConfig)
	topSeat := helper.RandInt(bigOrSmallConfigLen)

	newBigOrSmallConfigs := make([]*BigOrSmallConfig, 0)
	newBigOrSmallConfigs = append(newBigOrSmallConfigs, trs.BigOrSmallConfig[topSeat:]...)
	newBigOrSmallConfigs = append(newBigOrSmallConfigs, trs.BigOrSmallConfig[:topSeat]...)

	for k, _ := range newBigOrSmallConfigs {
		newBigOrSmallConfigs[k].Seat = k
	}

	//trs.ComRoomSpace.CurrAnimalWheelSort = append(trs.ComRoomSpace.CurrAnimalWheelSort, newAnimalConfigs)
	return newBigOrSmallConfigs
}

func WheelAnimalPartyCalculateExec(trs *RoomSpace) {
	var (
		//赢钱的用户的 ｜ 输钱的用户
		//放在一起返回
		currPeriodUserWinMsg = &pbs.CurrPeriodUserWinMsg{UserBetSettle: make([]*pbs.UserBetSettle, 0)}

		//用户｜金额
		userWinLoseInfo = make(map[string]float32)

		//全部用户的押注情况 (每个用户的多个区域)
		allUserBetInfo = make(map[string]float32)

		//todo
		//每个区域的情况
	)

	for _, mapUInfo := range trs.ComRoomSpace.TurnMateInfo.BetZoneUserInfoMap {
		for _, uInfo := range mapUInfo {
			currVal, ok := allUserBetInfo[uInfo.UserID]
			if ok {
				allUserBetInfo[uInfo.UserID] = float32(uInfo.UserProperty.Bet) + currVal
			} else {
				allUserBetInfo[uInfo.UserID] = float32(uInfo.UserProperty.Bet)
			}
		}
	}

	//获取所有的用户押注情况
	winBetZoneConfig := trs.ComRoomSpace.CurrAnimalWheelSort

	//动物赢钱的区域
	for k, animalWheelSort := range winBetZoneConfig {
		if k == 0 {
			// 1=大（粉色） 2=小（紫色）
			if animalWheelSort.WinBigOrSmallConfig.BigOrSmall == 1 {
				//8
				winUserArr, _ := trs.ComRoomSpace.GetBetZoneUserInfos(8)
				//押注大小的赢钱区域
				for _, uInfo := range winUserArr {
					currVal, ok := userWinLoseInfo[uInfo.UserID]
					if ok {
						userWinLoseInfo[uInfo.UserID] = float32(helper.Mul(uInfo.UserProperty.Bet, 2)) + currVal
					} else {
						userWinLoseInfo[uInfo.UserID] = float32(helper.Mul(uInfo.UserProperty.Bet, 2))
					}
				}

				//for _, uInfo := range loseUserArr {
				//	currVal, ok := userWinLoseInfo[uInfo.UserID]
				//	if ok {
				//		userWinLoseInfo[uInfo.UserID] = float32(helper.Sum(-float32(uInfo.UserProperty.Bet), currVal))
				//	} else {
				//		userWinLoseInfo[uInfo.UserID] = -float32(uInfo.UserProperty.Bet)
				//	}
				//}
			}
			if animalWheelSort.WinBigOrSmallConfig.BigOrSmall == 2 {
				//12
				winUserArr, _ := trs.ComRoomSpace.GetBetZoneUserInfos(12)
				//押注大小的赢钱区域
				for _, uInfo := range winUserArr {
					currVal, ok := userWinLoseInfo[uInfo.UserID]
					if ok {
						userWinLoseInfo[uInfo.UserID] = float32(helper.Mul(uInfo.UserProperty.Bet, 2)) + currVal
					} else {
						userWinLoseInfo[uInfo.UserID] = float32(helper.Mul(uInfo.UserProperty.Bet, 2))
					}
				}

				//for _, uInfo := range loseUserArr {
				//	currVal, ok := userWinLoseInfo[uInfo.UserID]
				//	if ok {
				//		userWinLoseInfo[uInfo.UserID] = float32(helper.Sum(-float32(uInfo.UserProperty.Bet), currVal))
				//	} else {
				//		userWinLoseInfo[uInfo.UserID] = -float32(uInfo.UserProperty.Bet)
				//	}
				//}

			}

		}

		if animalWheelSort.WinAnimalConfig.AnimalId == 2 {
			//幸运动物
			continue
		}
		//常规动物的赢钱
		winBetZoneConf := animalWheelSort.WinBetZoneConfig
		for _, winBetZone := range winBetZoneConf {
			//先发中奖组合
			winUserArr, _ := trs.ComRoomSpace.GetBetZoneUserInfos(winBetZone.Seat)
			//global.GVA_LOG.Infof("WheelAnimalPartyCalculateExec 中奖用户 winUserArr: %v , loseUserAr %v", winUserArr, loseUserArr)

			for _, uInfo := range winUserArr {
				currVal, ok := userWinLoseInfo[uInfo.UserID]
				if ok {
					userWinLoseInfo[uInfo.UserID] = float32(helper.Mul(uInfo.UserProperty.Bet, winBetZone.BetRate)) + currVal
				} else {
					userWinLoseInfo[uInfo.UserID] = float32(helper.Mul(uInfo.UserProperty.Bet, winBetZone.BetRate))
				}
			}

			//for _, uInfo := range loseUserArr {
			//	currVal, ok := userWinLoseInfo[uInfo.UserID]
			//	if ok {
			//		userWinLoseInfo[uInfo.UserID] = float32(helper.Sum(-float32(uInfo.UserProperty.Bet), currVal))
			//	} else {
			//		userWinLoseInfo[uInfo.UserID] = -float32(uInfo.UserProperty.Bet)
			//	}
			//}
		}
	}

	for uId, betAll := range allUserBetInfo {
		winVal, ok := allUserBetInfo[uId]
		if ok {
			if winVal > betAll {
				currPeriodUserWinMsg.UserBetSettle = append(currPeriodUserWinMsg.UserBetSettle, &pbs.UserBetSettle{
					WinCoin: winVal - betAll,
					UserId:  uId,
				})
			} else {
				currPeriodUserWinMsg.UserBetSettle = append(currPeriodUserWinMsg.UserBetSettle, &pbs.UserBetSettle{
					LoseCoin: winVal - betAll,
					UserId:   uId,
				})
			}
		} else {
			currPeriodUserWinMsg.UserBetSettle = append(currPeriodUserWinMsg.UserBetSettle, &pbs.UserBetSettle{
				LoseCoin: -betAll,
				UserId:   uId,
			})
		}

	}

	//for uid, coinNum := range userWinLoseInfo {
	//	if coinNum >= 0 {
	//		currPeriodUserWinMsg.UserBetSettle = append(currPeriodUserWinMsg.UserBetSettle, &pbs.UserBetSettle{
	//			WinCoin: coinNum,
	//			UserId:  uid,
	//		})
	//	} else {
	//		currPeriodUserWinMsg.UserBetSettle = append(currPeriodUserWinMsg.UserBetSettle, &pbs.UserBetSettle{
	//			LoseCoin: coinNum,
	//			UserId:   uid,
	//		})
	//	}
	//}

	netMessageResp := helper.NewNetMessage("", "", int32(pbs.ProtocNum_CurrPeriodUserWinMsg), config.SlotServer)
	responseHeadByte, _ := proto.Marshal(currPeriodUserWinMsg)
	netMessageResp.Content = responseHeadByte
	NatsSendAimUserMsg(trs, netMessageResp, "")

}

func (trs *RoomSpace) CalculateAndEnd() {
	//计算前面所有轮的点赞情况 得出排名
	var userLikeCardMap []*models.UserLikeDetail
	userLikeCardMap = make([]*models.UserLikeDetail, 0)

	//一个赞 50分
	var userIntegral map[string]float64
	userIntegral = make(map[string]float64)

	//连续点赞的纪录
	var isOnGoLike map[string][]bool
	isOnGoLike = make(map[string][]bool)

	//var userExperience map[string]float64
	//userExperience = make(map[string]float64)

	//1 获取前面几轮 被点赞的牌
	//allTurnLikeCards := trs.ComRoomSpace.AllTurnLikeCards()
	for i := 1; i <= trs.RoomInfo.RoomTurnNum; i++ {
		turnLikeCards := trs.ComRoomSpace.TurnLikeCards(i)
		for k, _ := range turnLikeCards {
			likeCard := turnLikeCards[k]
			if likeCard.LikeNum > 0 {
				integral, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", helper.Sum(userIntegral[likeCard.LikeUserId], helper.Mul(likeCard.LikeNum, 50))), 64)
				//去掉加成
				//if likeCard.Level == 2 {
				//	integral, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", helper.Mul(integral, 1.25)), 64)
				//}
				//if likeCard.Level == 3 {
				//	integral, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", helper.Mul(integral, 1.5)), 64)
				//}
				userIntegral[likeCard.LikeUserId] = integral
				isOnGoLike[likeCard.LikeUserId] = append(isOnGoLike[likeCard.LikeUserId], true)
			} else {
				isOnGoLike[likeCard.LikeUserId] = append(isOnGoLike[likeCard.LikeUserId], false)
			}
		}
	}

	//2 统计结果
	for _, uInfo := range trs.ComRoomSpace.UserInfos {
		fen := userIntegral[uInfo.UserID]
		boolArr := isOnGoLike[uInfo.UserID]
		maxLen := longestTrue(boolArr)
		resLikeDetail := models.UserLikeDetail{
			UserID:      uInfo.UserID,
			Nickname:    uInfo.Nickname,
			HeadPhoto:   uInfo.UserID,
			OnGoLinkNum: maxLen,
			Integral:    fen,
			Experience:  0,
			MCoin:       0,
		}

		userLikeCardMap = append(userLikeCardMap, &resLikeDetail)
	}

	//根据积分排序
	//降序排序
	sort.Slice(userLikeCardMap, func(i, j int) bool {
		return userLikeCardMap[i].Integral > userLikeCardMap[j].Integral
	})

	//根基排名 获取经验值个和积分
	calculateRank(userLikeCardMap)

	for k, _ := range userLikeCardMap {
		val := userLikeCardMap[k]
		info, ok := trs.ComRoomSpace.UserInfos[val.UserID]
		if !ok {
			continue
		}
		if info.UserIsRobot() {
			continue
		}
		err := dao.UpdateUserCoinNumOrExperience(val.UserID, val.MCoin, val.Experience, 0)
		if err != nil {
			global.GVA_LOG.Error("CalculateAndEnd ,UpdateUserCoinNumOrExperience err", zap.Any("err", err))
		}
	}

	netMessageResp := helper.NewNetMessage("", "", int32(pbs.Meb_calculateRank), config.SlotServer)
	//发送广播
	msgData := models.CalculateRankMsg{
		ProtoNum:       strconv.Itoa(int(pbs.Meb_calculateRank)),
		RoomNo:         trs.RoomInfo.RoomNo,
		Timestamp:      time.Now().Unix(),
		LikeDetailList: userLikeCardMap,
	}

	//给用户消息
	global.GVA_LOG.Infof("EnCalculateAndEnd-本轮最终用户排行计算: %v", msgData)
	responseHeadByte, _ := json.Marshal(msgData)
	netMessageResp.Content = responseHeadByte
	NatsSendAllUserMsg(trs, netMessageResp)
}

func calculateRank(userLikeCardMap []*models.UserLikeDetail) {
	//每局奖励：
	//第一：金币×100 + 经验×150
	//第二：金币×70 + 经验×100
	//第三：金币×40 + 经验×60
	//第四：金币×20 + 经验×30
	//
	//1-10级玩家：第4名奖励提升至第3名水平，降低前期流失率。 todo
	//完全平局（四人同排名）：所有玩家获得 第2名保底奖励（金币×70，经验×100）。
	//并列第1名：按原第1名奖励全额发放，不均分 每人获得：金币×100（非均分50），
	//后续名次跳过被占用的名次，剩余玩家按实际顺位结算 原第3名变为第2名，但奖励按第3名标准发放后续依此类推。

	//第一名获取分数
	const firstExperience = 150
	const secondExperience = 100
	const thirdExperience = 60
	const fourExperience = 30

	const firstCoin = 100
	const secondCoin = 70
	const thirdCoin = 40
	const fourCoin = 20

	//第一名得分 = 第二名得分
	firstExperienceQqSecond := false
	//第一名得分 = 第二名得分 = 第三名得分
	firstExperienceQqThird := false
	//第一名得分 = 第二名得分 = 第三名得分 = 第四名
	firstExperienceQqFour := false

	//第一名得分 != 第二名得分 = 第三名得分
	secondExperienceQqThird := false

	for k, _ := range userLikeCardMap {
		val := userLikeCardMap[k]
		if k == 0 {
			val.Experience = firstExperience
			val.MCoin = firstCoin
		}
		//排名第二的用户
		if k == 1 {
			//检测是否和第一名得分一样
			integral0 := userLikeCardMap[0].Integral
			if val.Integral == integral0 {
				val.Experience = userLikeCardMap[0].Experience
				val.MCoin = userLikeCardMap[0].MCoin
				firstExperienceQqSecond = true
			} else {
				val.Experience = secondExperience
				val.MCoin = secondCoin
			}
		}

		//排名第三的用户
		// 1 是否与第一名得分一样
		// 2 是否与第二名得分一样
		if k == 2 {
			integral0 := userLikeCardMap[0].Integral
			if firstExperienceQqSecond {
				//前2名得分一样
				if val.Integral == integral0 {
					val.Experience = userLikeCardMap[0].Experience
					val.MCoin = userLikeCardMap[0].MCoin
					firstExperienceQqThird = true
					continue
				}
			}

			integral1 := userLikeCardMap[1].Integral
			//第二名是否与第三名得分一样
			if integral1 == val.Integral {
				val.Experience = userLikeCardMap[1].Experience
				val.MCoin = userLikeCardMap[1].MCoin
				secondExperienceQqThird = true
			} else {
				val.Experience = thirdExperience
				val.MCoin = thirdCoin
			}
		}

		//排名第四的用户
		//1 是否与第一名得分一样
		//2 是否与第二名得分一样
		//3 是否与第三名得分一样
		if k == 3 {
			integral0 := userLikeCardMap[0].Integral
			//是否与第一名得分一样
			if firstExperienceQqThird {
				if val.Integral == integral0 {
					val.Experience = userLikeCardMap[0].Experience
					val.MCoin = userLikeCardMap[0].MCoin
					firstExperienceQqFour = true
					continue
				}
			}

			integral1 := userLikeCardMap[1].Integral
			//是否与第二名得分一样
			if secondExperienceQqThird {
				//第三名 是否 等于第四名
				if integral1 == val.Integral {
					val.Experience = userLikeCardMap[1].Experience
					val.MCoin = userLikeCardMap[1].MCoin
					secondExperienceQqThird = true
				} else {
					val.Experience = fourExperience
					val.MCoin = fourCoin
				}
				continue
			}

			//是否与第三名得分一样
			integral2 := userLikeCardMap[2].Integral
			if integral2 == val.Integral {
				val.Experience = userLikeCardMap[2].Experience
				val.MCoin = userLikeCardMap[2].MCoin
				secondExperienceQqThird = true
			} else {
				val.Experience = fourExperience
				val.MCoin = fourCoin
			}
		}
	}

	if firstExperienceQqFour == true {
		global.GVA_LOG.Infof("全部用户得分一样")
		for k, _ := range userLikeCardMap {
			val := userLikeCardMap[k]
			val.Experience = secondExperience
			val.MCoin = secondCoin
		}
	}

}

func longestTrue(slice []bool) int {
	maxLen, current := 0, 0
	for _, v := range slice {
		if v {
			current++
			if current > maxLen {
				maxLen = current
			}
		} else {
			current = 0
		}
	}
	return maxLen
}

//==

// EntryLikePage 是否该发送进入点赞页面的广播
//func (trs *RoomSpace) EntryLikePage() {
//	//if trs.ComRoomSpace.GetGameState() != EnLikePage {
//	//	return
//	//}
//	global.GVA_LOG.Infof("EnLikePage 房间 {%v} 第{%v} 轮 都出过牌了 现在发送广播 进入点赞页面", trs.RoomInfo.RoomNo, trs.ComRoomSpace.GetTurn())
//
//	//获取本轮房间所有用户的牌，发送给房间所有人
//	allOutCart := []*models.Card{}
//
//	for _, uInfo := range trs.ComRoomSpace.UserInfos {
//		cards := trs.ComRoomSpace.GetUserOutEdCards(uInfo.UserID)
//		for _, card := range cards {
//			card.UserID = uInfo.UserID
//			allOutCart = append(allOutCart, card)
//		}
//	}
//	//发送广播
//	msgData := models.EntryLikePageMsg{
//		ProtoNum:  strconv.Itoa(int(pbs.Meb_entryLikePage)),
//		RoomNo:    trs.RoomInfo.RoomNo,
//		Timestamp: time.Now().Unix(),
//		OutCards:  allOutCart,
//	}
//	//给用户消息
//	global.GVA_LOG.Infof("StartPlay 开始游戏的广播: %v", msgData)
//	responseHeadByte, _ := json.Marshal(msgData)
//	NatsSendAllUserMsg(trs, helper.GetNetMessage(0, 0, int32(pbs.Meb_entryLikePage), config.NatsMemeBattle, responseHeadByte))
//
//	trs.ComRoomSpace.GameStateTransition(EnLikePage, DisLikePage)
//}
//
//func (trs *RoomSpace) NextTurnOrCalculateAndEnd() {
//	//if trs.ComRoomSpace.GetGameState() != EnLikeCard {
//	//	return
//	//}
//	global.GVA_LOG.Infof("NextTurnOrCalculateAndEnd 房间 {%v} 第{%v} 轮 都点过赞了", trs.RoomInfo.RoomNo, trs.ComRoomSpace.GetTurn())
//
//	if trs.RoomInfo.RoomTurnNum == trs.ComRoomSpace.GetTurn() {
//		//最后一轮进入结算状态
//		trs.ComRoomSpace.GameStateTransition(EnLikeCard, EnCalculateAndEnd)
//
//		//已经是最后一轮就是结束
//		trs.CalculateAndEnd()
//
//		//最近一轮的点赞结束 触发游戏结束
//		if trs.RoomInfo.RoomTurnNum == trs.ComRoomSpace.GetTurn() {
//			trs.CloseRoom(trs.RoomInfo.RoomNo, table.RoomStatusStop)
//		}
//	} else {
//		//本轮点赞阶段结束
//		trs.ComRoomSpace.GameStateTransition(EnLikeCard, DisLikeCard)
//
//		//如果还没到达最后一轮 就是进入下一轮
//		//发送问题
//		trs.NextTurn()
//
//		//房间维度
//		//问题阶段，收到问题 还没收到牌
//		trs.ComRoomSpace.GameStateTransition(DisLikeCard, IssueStage)
//	}
//
//}

//func NextTurnOrCalculateAndEnd(trs *RoomSpace) {
//	global.GVA_LOG.Infof("NextTurnOrCalculateAndEnd 房间 {%v} 第{%v} 轮 都点过赞了", trs.RoomInfo.RoomNo, trs.ComRoomSpace.GetTurn())
//
//	if trs.RoomInfo.RoomTurnNum == trs.ComRoomSpace.GetTurn() {
//		//最后一轮进入结算状态
//		trs.ComRoomSpace.GameStateTransition(EnLikeCard, EnCalculateAndEnd)
//
//		//已经是最后一轮就是结束
//		trs.CalculateAndEnd()
//
//		//最近一轮的点赞结束 触发游戏结束
//		if trs.RoomInfo.RoomTurnNum == trs.ComRoomSpace.GetTurn() {
//			trs.CloseRoom(trs.RoomInfo.RoomNo, table.RoomStatusStop)
//		}
//	} else {
//		//本轮点赞阶段结束
//		trs.ComRoomSpace.GameStateTransition(EnLikeCard, DisLikeCard)
//
//		//如果还没到达最后一轮 就是进入下一轮
//		//发送问题
//		trs.NextTurn()
//
//		//房间维度
//		//问题阶段，收到问题 还没收到牌
//		trs.ComRoomSpace.GameStateTransition(DisLikeCard, IssueStage)
//	}
//}

//func (trs *RoomSpace) NextTurn() {
//	//房间轮增加
//	trs.ComRoomSpace.AddTurn()
//	trs.ComRoomSpace.SetCountdownTime(helper.LocalTime().Unix())
//
//	////用户出牌倒计时
//	//for _, uInfo := range trs.ComRoomSpace.UserInfos {
//	//	uInfo.SetOutCardCountDown(models.GetOutCardCountDownTimeInt(OutCardCountDownTimeInt))
//	//}
//
//	//发问题 ｜ 发牌逻辑在定时任务里面
//	issue, err := trs.SelectIssue()
//	if err != nil {
//		global.GVA_LOG.Infof("NextTurn %v", zap.Any("err", err))
//	}
//
//	global.GVA_LOG.Infof("NextTurn %v", zap.Any("issue", &issue))
//}

//func WheelAnimalPartyCalculateExec(trs *RoomSpace) {
//	//获取所有的用户押注情况
//	winBetZoneConfig := trs.ComRoomSpace.WinBetZoneConfig
//
//	//先发中奖组合
//	winUserArr, loseUserArr := trs.ComRoomSpace.GetBetZoneUserInfos(winBetZoneConfig.Seat)
//	global.GVA_LOG.Infof("WheelAnimalPartyCalculateExec 中奖用户 winUserArr: %v , loseUserAr %v", winUserArr, loseUserArr)
//
//	currPeriodUserWinMsg := &pbs.CurrPeriodUserWinMsg{
//		UserBetSettle: make([]*pbs.UserBetSettle, 0),
//	}
//
//	for _, uInfo := range winUserArr {
//		currPeriodUserWinMsg.UserBetSettle = append(currPeriodUserWinMsg.UserBetSettle, &pbs.UserBetSettle{
//			WinCoin: float32(helper.Sum(uInfo.UserProperty.Bet, winBetZoneConfig.BetRate)),
//			UserId:  uInfo.UserID,
//		})
//	}
//
//	for _, uInfo := range loseUserArr {
//		currPeriodUserWinMsg.UserBetSettle = append(currPeriodUserWinMsg.UserBetSettle, &pbs.UserBetSettle{
//			LoseCoin: float32(uInfo.UserProperty.Bet),
//			UserId:   uInfo.UserID,
//		})
//	}
//
//	netMessageResp := helper.NewNetMessage("", "", int32(pbs.ProtocNum_CurrPeriodUserWinMsg), config.SlotServer)
//
//	responseHeadByte, _ := proto.Marshal(currPeriodUserWinMsg)
//	netMessageResp.Content = responseHeadByte
//	NatsSendAimUserMsg(trs, netMessageResp, "")
//
//}

//func AnimalPartyCalculateExecFunc(trs *RoomSpace) {
//	global.GVA_LOG.Infof("AnimalPartyCalculateExecFunc 房间 {%v}", trs.RoomInfo.RoomNo)
//
//	//计算逻辑
//
//	//当前指针所在方向
//
//	//当前所有用户的押注分布
//
//	//哪些用户赢钱
//
//	//保存数据库
//
//	//推送消息
//}

// CalculateExecFunc  计算 并结束
//func CalculateExecFunc(trs *RoomSpace) {
//	global.GVA_LOG.Infof("CalculateExecFunc 房间 {%v} 第{%v} 轮 都点过赞了", trs.RoomInfo.RoomNo, trs.ComRoomSpace.GetTurn())
//
//	//已经是最后一轮就是结束
//	trs.CalculateAndEnd()
//
//	//最近一轮的点赞结束 触发游戏结束
//	if trs.RoomInfo.RoomTurnNum == trs.ComRoomSpace.GetTurn() {
//		trs.CloseRoom(trs.RoomInfo.RoomNo, table.RoomStatusStop)
//	}
//
//}

// NextTurnExecFunc 进入下一轮
//func NextTurnExecFunc(trs *RoomSpace) {
//	global.GVA_LOG.Infof("NextTurnExecFunc 房间 {%v} 第{%v} 轮 都点过赞了", trs.RoomInfo.RoomNo, trs.ComRoomSpace.GetTurn())
//	//如果还没到达最后一轮 就是进入下一轮
//	//trs.NextTurn()
//
//	//房间轮增加
//	trs.ComRoomSpace.AddTurn()
//	trs.ComRoomSpace.SetCountdownTime(helper.LocalTime().Unix()) //进入下一轮
//	trs.ComRoomSpace.SetLikeCountdownTime(0)
//
//	//做2件事情
//	//1 发送问题
//	//2 发牌
//	trs.ExecSendIssueAndSendCards() //进入下一轮
//
//}

//func SendGameStartBroadcast(trs *RoomSpace) {
//	var (
//		userInfos         = trs.ComRoomSpace.UserInfos
//		memeRoomUserLists []models.MemeRoomUser
//		index             int
//	)
//
//	err := table.SaveMemeRoom(trs.RoomInfo)
//	if err != nil {
//		global.GVA_LOG.Error("StartPlay SaveTavernRoom ", zap.Error(err))
//	}
//
//	//先获取房间全部的用户
//	for k, _ := range userInfos {
//		userItem := userInfos[k]
//		userItem.UserProperty.Turn = trs.ComRoomSpace.GetTurn()
//		tavernRoomUser := models.MemeRoomUser{
//			UserID:       userItem.UserID,
//			Nickname:     userItem.Nickname,
//			Turn:         userItem.UserProperty.Turn,
//			IsLeave:      userItem.UserProperty.IsLeave,
//			IsOwner:      userItem.UserProperty.IsOwner,
//			IsReady:      userItem.UserProperty.IsReady,
//			Seat:         userItem.UserProperty.Seat,
//			UserLimitNum: userItem.UserProperty.UserLimitNum,
//			WinPrice:     userItem.UserProperty.WinPrice,
//			Bet:          userItem.UserProperty.Bet,
//		}
//		memeRoomUserLists = append(memeRoomUserLists, tavernRoomUser)
//		index++
//	}
//
//	//发送广播
//	msgData := models.StartPlayMsg{
//		ProtoNum:         strconv.Itoa(int(pbs.Meb_startPlay)),
//		RoomNo:           trs.RoomInfo.RoomNo,
//		Timestamp:        time.Now().Unix(),
//		MemeRoomUserList: memeRoomUserLists,
//	}
//
//	//给用户消息
//	global.GVA_LOG.Infof("StartPlay 开始游戏的广播: %v", msgData)
//	responseHeadByte, _ := json.Marshal(msgData)
//	NatsSendAllUserMsg(trs, helper.GetNetMessage("", "", int32(pbs.Meb_startPlay), config.SlotServer, responseHeadByte))
//
//	trs.ComRoomSpace.SetGameStartTime(helper.LocalTime().Unix()) //游戏开始时间
//}

//func SendEnLoadBroadcast(trs *RoomSpace) {
//	//如果都加载完成 需要通知客户端 并发送问题
//	//trs.AllUserLoad()
//	trs.ComRoomSpace.SetCountdownTime(helper.LocalTime().Unix()) //加载完成
//
//	//做2件事情
//	//1 发送问题
//	//2 发牌
//	//trs.ExecSendIssueAndSendCards() //加载完成
//}

//func (trs *RoomSpace) ExecSendIssueAndSendCards() {
//	var (
//		roomNo = trs.RoomInfo.RoomNo
//		turn   = trs.ComRoomSpace.GetTurn()
//	)
//
//	//trs.SaveAndSendIssue()
//	//
//	//trs.InitUserSelfCards()
//
//	netMessageResp := helper.NewNetMessage("", "", int32(pbs.Meb_dealCardsMsg), config.SlotServer)
//	//发送广播
//	for _, uInfo := range trs.ComRoomSpace.UserInfos {
//		//获取当前轮 用户 没有被随的牌
//		cards, err := trs.ComRoomSpace.GetNotExtractCard(uInfo.UserID)
//		if err != nil {
//			global.GVA_LOG.Infof("DealCards %v", err.Error())
//			continue
//		}
//		if cards != nil && len(cards) <= 0 {
//			continue
//		}
//
//		//所有牌随机打乱
//		helper.SliceShuffle(cards)
//		var (
//			newCards = make([]*table.MbCardConfig, 0) //未被选的牌
//			outCards = make([]*table.MbCardConfig, 0) //给用户要发的牌
//			perCards = make([]*models.Card, 0)        //每个人获取4张牌
//		)
//
//		for _, val := range cards {
//			if len(perCards) < 4 {
//				v := &models.Card{
//					CardId:  val.ID,
//					Name:    val.Name,
//					Suffix:  val.SuffixName,
//					Level:   val.Level,
//					AddRate: val.AddRate,
//					UserID:  uInfo.UserID,
//				}
//				perCards = append(perCards, v)
//				outCards = append(outCards, val)
//			} else {
//				newCards = append(newCards, val)
//			}
//		}
//
//		err = trs.ComRoomSpace.AddCurrCard(uInfo.UserID, perCards)
//		if err != nil {
//			global.GVA_LOG.Error("dealCards  AddCurrCard", zap.Error(err))
//		}
//
//		//重置 未抽过的牌
//		trs.ComRoomSpace.ReMakeExtractCard(uInfo.UserID, newCards)
//
//		//保留抽过的牌
//		trs.ComRoomSpace.SaveExtractCard(uInfo.UserID, outCards)
//
//		//发牌:给每一个用户发对应的牌
//		msgData := models.DealCardsMsg{
//			ProtoNum:  strconv.Itoa(int(pbs.Meb_dealCardsMsg)),
//			Timestamp: time.Now().Unix(),
//			UserId:    uInfo.UserID,
//			RoomNo:    roomNo,
//			Turn:      turn,
//			Cards:     perCards,
//		}
//
//		global.GVA_LOG.Infof("发牌的广播: %v", msgData)
//		userStateRespMarshal, _ := json.Marshal(msgData)
//		netMessageResp.Content = userStateRespMarshal
//		NatsSendAimUserMsg(trs, netMessageResp, uInfo.UserID)
//
//		//发完牌纪录一下具体数据 mysql （暂时先不纪录）
//		//perCardsMarshal, _ := json.Marshal(perCards)
//		//dao.AddTurnDetails(roomNo, uInfo.UserID, uInfo.Nickname, turn, string(perCardsMarshal), "{}")
//		//trs.ByRobotClassSetAction(uInfo)
//	}
//}

//func (trs *RoomSpace) SaveAndSendIssue() {
//	//问题
//	issue, err := trs.ComRoomSpace.GetSelectIssue()
//	if err != nil {
//		global.GVA_LOG.Infof("SelectIssue  GetSelectIssue %v", zap.Error(err))
//	} else {
//		//返回已经存在的骗子牌
//		global.GVA_LOG.Error("SelectIssue 洗牌提前 GetFraudCard", zap.Error(err))
//		return
//	}
//
//	randInt := helper.RandInt(len(trs.RoomIssueConfig))
//	issueConfig := trs.RoomIssueConfig[randInt]
//	issue = &models.Issue{
//		IssueId: issueConfig.ID,
//		Level:   issueConfig.Level,
//		Class:   issueConfig.Class,
//		Desc:    issueConfig.Desc,
//	}
//
//	//保存问题
//	trs.ComRoomSpace.AddSelectIssue(issue)
//
//	msgData := models.IssueMsg{
//		ProtoNum:  strconv.Itoa(int(pbs.Meb_issueMsg)),
//		Timestamp: time.Now().Unix(),
//		Issue:     issue,
//	}
//	global.GVA_LOG.Infof("SelectIssue 本轮问题的的广播: %v ", msgData)
//	responseHeadByte, _ := json.Marshal(msgData)
//	NatsSendAllUserMsg(trs, helper.GetNetMessage("", "", int32(pbs.Meb_issueMsg), config.SlotServer, responseHeadByte)) //SelectIssue
//}

// InitUserSelfCards 每轮开始前 初始化自己的牌
//func (trs *RoomSpace) InitUserSelfCards() {
//	for _, uInfo := range trs.ComRoomSpace.UserInfos {
//		//先把基础牌放到未随机里面
//		notExtractCards := make([]*table.MbCardConfig, 0)
//		//基础牌 去掉了
//		//notExtractCards = append(notExtractCards, trs.RoomBaseCard...)
//
//		if uInfo.UserIsRobot() {
//			//获取当前的版本
//			for _, vCards := range trs.RoomVersionCard {
//				notExtractCards = append(notExtractCards, vCards...)
//			}
//		} else {
//			userCards := trs.ComRoomSpace.GetUserOwnCards(uInfo.UserID)
//			if len(userCards) == 0 {
//				//查找自己的牌并赋值
//				cardConfigByIds := logic.GetUserOwnCards(uInfo.UserID)
//				trs.ComRoomSpace.AddUserOwnCards(uInfo.UserID, cardConfigByIds)
//				userCards = cardConfigByIds
//			}
//			//自己的牌
//			notExtractCards = append(notExtractCards, userCards...)
//		}
//
//		err := trs.ComRoomSpace.CurrTurnFirstNotExtractCard(uInfo.UserID, notExtractCards)
//		if err != nil {
//			global.GVA_LOG.Infof("DealCards %v", err.Error())
//			continue
//		}
//	}
//}

// EntryLikePage 是否该发送进入点赞页面的广播
//func EntryLikePage(trs *RoomSpace) {
//	global.GVA_LOG.Infof("EnLikePage 房间 {%v} 第{%v} 轮 都出过牌了 现在发送广播 进入点赞页面", trs.RoomInfo.RoomNo, trs.ComRoomSpace.GetTurn())
//
//	//获取本轮房间所有用户的牌，发送给房间所有人
//	allOutCard := trs.ComRoomSpace.GetAllUserOutEdCards()
//
//	//发送广播
//	msgData := models.EntryLikePageMsg{
//		ProtoNum:  strconv.Itoa(int(pbs.Meb_entryLikePage)),
//		RoomNo:    trs.RoomInfo.RoomNo,
//		Timestamp: time.Now().Unix(),
//		OutCards:  allOutCard,
//	}
//	//给用户消息
//	global.GVA_LOG.Infof("StartPlay 开始游戏的广播: %v", msgData)
//	responseHeadByte, _ := json.Marshal(msgData)
//	NatsSendAllUserMsg(trs, helper.GetNetMessage("", "", int32(pbs.Meb_entryLikePage), config.SlotServer, responseHeadByte))
//
//	//点赞倒计时
//	trs.ComRoomSpace.SetLikeCountdownTime(helper.LocalTime().Unix())
//
//	//通知完成
//	//trs.ComRoomSpace.GameStateTransition(EnLikePage, DisLikePage)
//	////并进入点赞中
//	//trs.ComRoomSpace.GameStateTransition(DisLikePage, EnLikeCardIng)
//}
