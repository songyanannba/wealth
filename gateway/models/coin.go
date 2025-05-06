package models

const TavernStoryQuickMatch = "31" //快速匹配 减积分
const TavernStoryRoomMatch = "32"  //房间匹配 减积分
const TavernStoryReturnCoin = "33" //返还积分（游戏没开始前 离开房间） 加积分
const TavernStoryWinAddCoin = "35" //最后赢家 加积分

type UpdateUserScore struct {
	Uid    string `json:"uid"`     //
	RoomNo string `json:"room_no"` //房间编号
	Bet    string `json:"bet"`     //押注
	Settle string `json:"settle"`  //
	//游戏过程:
	//|| (1-4:大逃杀) 1=投注,2=主动退出,3=吃鸡 4=主动增加积分
	//||（11-20：钓鱼） 11 挥杆消耗 12 挥杆结果收益（如果钓的鱼很值钱，这边是大于挥杆消耗的，如果不值钱小于挥杆消耗）13 开宝箱收益 14 福鱼
	//|| (21-30：挖矿）21:金矿销毁补偿 （矿石兑换金币）；22:金矿购买 减自己的金币（转增）; 23:金矿购买 加售卖人的金币（转增）; 24.挖金矿矿工购买（购买矿工） 25:金矿游戏补偿积分 26:重置年龄
	//|| (31-40：骗子酒店）31:快速匹配模式：扣减积分  32:房间模式：开始游戏扣减积分 33：返还用户积分 35:游戏结束赢家增加积分
	Process string `json:"process"` //
	Nemo    string `json:"nemo"`    //
}
