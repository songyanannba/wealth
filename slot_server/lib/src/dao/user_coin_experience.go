package dao

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models/table"
)

func UpdateUserCoinNum(uid string, coin float64) error {
	global.GVA_LOG.Infof("UpdateUserCoinNum uid:%v ,coin:%v", uid, coin)

	updateMap := MakeUpdateData("coin_num", coin)

	err := table.UpdateUserCoinExperience(uid, updateMap)
	if err != nil {
		global.GVA_LOG.Error("UpdateUserCoinNum ", zap.Error(err))
		return err
	}
	return nil
}

func UpdateUserExperience(uid string, experience float64) error {
	global.GVA_LOG.Infof("UpdateUserExperience uid:%v  experience:%v", uid, experience)

	updateMap := MakeUpdateData("experience", experience)

	err := table.UpdateUserCoinExperience(uid, updateMap)
	if err != nil {
		global.GVA_LOG.Error("UpdateUserExperience ", zap.Error(err))
		return err
	}
	return nil
}

// UpdateUserCoinNumOrExperience 加分和经验 opType 1=开包 2=重随 0 加经验和积分
func UpdateUserCoinNumOrExperience(uid string, coinNum, experience float64, opType int) error {
	var (
		currCoinNum float64
		afterNum    float64
		beforeNum   float64
	)
	currCoinNum = coinNum
	userCoinExperience, err := table.GetUserCoinExperience(uid)
	if err != nil {
		global.GVA_LOG.Error("UpdateUserCoinNumOrExperience ", zap.Error(err))
	}
	beforeNum = userCoinExperience.CoinNum

	if userCoinExperience.ID > 0 {
		userCoinExperience.CoinNum = helper.Sum(userCoinExperience.CoinNum, coinNum)
		userCoinExperience.Experience = helper.Sum(userCoinExperience.Experience, experience)
		err := table.SaveUserCoinExperience(userCoinExperience)
		if err != nil {
			global.GVA_LOG.Error("userCoinExperience ", zap.Error(err))
		}
		afterNum = userCoinExperience.CoinNum
	} else {
		val := &table.UserCoinExperience{
			UserId:     uid,
			CoinNum:    coinNum,
			Experience: experience,
			DateTime:   helper.LocalTime(),
		}
		err := table.CreateUserCoinExperience(val)
		if err != nil {
			global.GVA_LOG.Error("UpdateUserCoinNumOrExperience tavernUserStatus", zap.Error(err))
		}
		afterNum = val.CoinNum
	}

	CreateCoinConsumeLog(uid, opType, currCoinNum, afterNum, beforeNum, "")
	return nil
}

func GetUserCoinExperience(uid string) *table.UserCoinExperience {
	userCoinExperience, err := table.GetUserCoinExperience(uid)
	if err != nil {
		global.GVA_LOG.Error("GetUserCoinExperience ", zap.Error(err))
	}
	return userCoinExperience
}
