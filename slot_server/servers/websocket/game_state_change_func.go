package websocket

import (
	"github.com/golang/protobuf/proto"
	"slot_server/lib/config"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/protoc/pbs"
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

	//winSeat = 16 //todo 测试

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

	//topSeat = 0 //todo 测试

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
		winVal, ok := userWinLoseInfo[uId]
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

	netMessageResp := helper.NewNetMessage("", "", int32(pbs.ProtocNum_CurrPeriodUserWinMsg), config.SlotServer)
	responseHeadByte, _ := proto.Marshal(currPeriodUserWinMsg)
	netMessageResp.Content = responseHeadByte
	NatsSendAimUserMsg(trs, netMessageResp, "")

}
