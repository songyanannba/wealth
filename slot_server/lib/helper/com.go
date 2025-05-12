package helper

import (
	"math"
	"slot_server/lib/common"
	"slot_server/protoc/pbs"
	"strconv"
)

// AttenuationByAge 根据年龄的衰减算法
func AttenuationByAge(age int) float64 {
	//baseValue := 1000.0
	//decayRate := 0.001

	//初始值为每小时5 刷机速度为8%
	//最低值为0.01

	baseValue := 5.0
	decayRate := 0.08
	age = age - 18
	// 计算衰减值
	decay := baseValue * math.Pow(1-decayRate, float64(age))
	if decay < 0.01 {
		decay = 0.01
	}
	return decay
}

func NewAttenuationByAge(age int, baseValue, decayRate float64) float64 {
	//初始值为每小时5 刷机速度为8%
	//最低值为0.01
	age = age - 18
	// 计算衰减值
	decay := baseValue * math.Pow(1-decayRate, float64(age))
	if decay < 0.01 {
		decay = 0.01
	}
	return decay
}

func NewNetMessage(reqUid, ackUid, msgId int32, serviceId string) *pbs.NetMessage {
	ackMsg := common.GetErrorMessage(common.OK, "")
	netMessageResp := &pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      reqUid,
			Token:    "",
			Platform: "",
		},
		AckHead: &pbs.AckHead{
			Uid:     ackUid,
			Code:    pbs.Code_OK,
			Message: ackMsg,
		},
		ServiceId: serviceId,
		MsgId:     msgId,
		Content:   nil,
	}
	return netMessageResp
}

func GetNetMessage(reqUid, ackUid, msgId int32, serviceId string, content []byte) *pbs.NetMessage {
	ackMsg := common.GetErrorMessage(common.OK, "")
	netMessageResp := &pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      reqUid,
			Token:    "",
			Platform: "",
		},
		AckHead: &pbs.AckHead{
			Uid:     ackUid,
			Code:    pbs.Code_OK,
			Message: ackMsg,
		},
		ServiceId: serviceId,
		MsgId:     msgId,
		Content:   content,
	}
	return netMessageResp
}

// GetIntUserId 转换用户id类型
func GetIntUserId(userId string) int {
	uIdInt, _ := strconv.Atoi(userId)
	return uIdInt
}
