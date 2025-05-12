// Package common 通用函数
package common

const (
	WebOK                       = 1    // Success 平台历史原因 原生返回的正确码都是 1
	OK                          = 200  // Success
	NotLoggedIn                 = 1000 // 未登录
	ParameterIllegal            = 1001 // 参数不合法
	UnauthorizedUserID          = 1002 // 非法的用户 ID
	Unauthorized                = 1003 // 未授权
	ServerError                 = 1004 // 系统错误
	NotData                     = 1005 // 没有数据
	ModelAddError               = 1006 // 添加错误
	ModelDeleteError            = 1007 // 删除错误
	ModelStoreError             = 1008 // 存储错误
	OperationFailure            = 1009 // 操作失败
	RoutingNotExist             = 1010 // 路由不存在
	RepetitiveOperation         = 1011 // 重复操作
	SysBusy                     = 1012 // 系统繁忙 稍后在试
	Maintenance                 = 1013 // 维护阶段 稍后再来
	UserScoreNotEnough          = 1014 // 用户积分不够
	SelectRoomProhibit          = 1015 // 房间被占用
	BetLow                      = 1017 // 押注不正确
	NoReSelectRoom              = 1018 // 已经选择过房间
	ReLogin                     = 1020
	LoginFirstProtectionDesc    = 1021 //首刀保护的档位 不让重复进入文案
	TokenExpiration             = 1022
	ServerRedisError            = 1023
	GetScoreErr                 = 1024 //获取远程服务错误
	DuplicateRequests           = 1025 //请求频繁
	GetRoomConfigErr            = 1026
	GetCurrTurnErr              = 1027 //获取当前轮错误
	GetTurnDetail               = 1028 //获取当前轮房间详情错误
	ProhibitSelectRoom          = 1029 // 系统强制匹配时间
	EnterGameIng                = 1030 //繁忙,请重试
	NotScoreConsume             = 1031 //耐久度不足
	RepetitionConfirm           = 1032 //重复确认
	DBErr                       = 1033 // 没有数据
	NoTFishRod                  = 1034 //
	LCoin                       = 1035 // 里昂尼斯币不足
	SCoin                       = 1036 // 国王积分不足
	CoinTypeNotMatchFishRodType = 1037 // 币类型和鱼竿类型不匹配
	NotSupportCoin              = 1038 // 不支持的币种
	ParameterNot                = 1039 // 参数有误
	FishRodNumNotEnough         = 1040 // 耐久度不够,请购买
	RepetitionFunc              = 1041 //重复操作
	SignInEd                    = 1042 //
	GetDBDataErr                = 1043 // 获取数据库数据有误
	FishSellTypeErr             = 1044 // 此类型鱼,不能出售
	UserKnapsackDataEmpty       = 1045 // 背包数据不足
	UserKnapsackCapacityLimit   = 1046 // 容量限制达到最大限制
	UserKnapsackStackUp         = 1047 // 背包叠放达到最大限制
	RpcCallRespErr              = 1048 // 用rpc调用远程服务器错误
	RpcCallRespDataErr          = 1049 // 用rpc调用远程服务器返回数据有误
	InvitesNumEnough            = 1050 // 邀请的用户数不够3人
	DbAddUserRatErr             = 1051 //
	RoomLevelNotEnough          = 1052 // 当前你能容纳的矿工数量已满，请升级您的金矿
	ProtocNumberError           = 1053 // 内部协议号错误
	DataCompileError            = 1054 // 数据编译错误
	GetDataFromDbErr            = 1055 // 获取数据错误
	AddDataFromDbErr            = 1056 // 添加数据错误
	RegisterTimeErr             = 1057 // 挖矿游戏上线之前已经注册
	KingCoinNotEnough           = 1058 // 没有足够金矿石
	UnauthorizedUserToken       = 1059 // 非法的用户token
	HaveBeInviteUser            = 1060 // 已经填写过邀请码
	InviteUserUnreal            = 1061 // 邀请用户不存在
	NotAuthentication           = 1062 // 请先去实名认证
	SellOreNumNot0              = 1063 // 售卖矿石数量必须大于0
	PrevLoginTimeTimeErr        = 1064 // 当前账号不是回归账号/新账号
	NotInviteYourself           = 1065 // 不允许邀请自己
	ExchangeOreCoinErr          = 1066 // 请输入正确的金矿数量
	BuyRatNumErr                = 1067 // 请输入正确的矿工数量
	TavernCreateRoomErr         = 1068 // 创建房间错误
	NotLogin                    = 1069 // 未登陆
	TavernRoomAlreadyFull       = 1070 // 已满员
	SellOreSinglePriceErr       = 1071 // 单价不能超过9999
	JoinRoomFull                = 1072 //  房间人数已满
	UserNotInRoom               = 1073 //  用户不在房间
	NotCanStartPlay             = 1074 //  不能开始游戏
	NotFraudCard                = 1075 //  还没有设置骗子牌
	ExistFraudCard              = 1076 //  已经存在骗子牌
	NotOutCartPower             = 1077 //  没有出牌的权限
	NotCanOutCards              = 1078 //  没有要出的牌
	NotDoubtCard                = 1079 //  没有质疑的权限
	ForceDoubt                  = 1080 //  上家牌出完,必须强制质疑
	AdoptNumErr                 = 1081 //  不能领取，收集的鱼不全
	BeOpenAdopt                 = 1082 //  已经开过宝箱
	BuyRatTimeLimit             = 1083 //  敬请期待 购买矿工的限制
	InviteGainRatLimit          = 1084 //  新用户需要购买一次矿工才能邀请成功
	NoWandErr                   = 1085 //  没有请求需要的标识
	NoWandPast                  = 1086 //  需要的标识过期
	CoinExchangeLimit           = 1087 //  每天积分总量限制
	CoinExchangePerKing         = 1088 //  每天积分达到限制
	CoinExchangePerOre          = 1089 //  每天矿石达到限制
	CoinBuySellLimit            = 1091 // 当天买卖总量限制
	CoinBuySellPerKing          = 1092 //  每天买卖积分达到限制
	CoinBuySellPerOre           = 1093 //  每天买卖矿石达到限制
	SellOreSinglePriceNotLess0  = 1094 //  单价不能小于0
	NOtOrder                    = 1095 //  订单不存在
	NOtOrderDone                = 1096 //  订单交易完成
	OreOrderStateDown           = 1097 //  挖矿售卖订单，上下架间隔最少5分钟
	ExchangeOreCoinNot9999      = 1098 //  每次兑换数量不能超过9999
	ServerTaskIng               = 1099 //  后台有任务在执行中 加入mysql执行队列，脚本执行
	NotSelfKickSelf             = 1100 // 自己不能踢自己
	NotRoomOwnKick              = 1101 // 不是房主 不能踢人
	RoomStatusStopErr           = 1102 // 房间游戏结束
	LeavePreRoom                = 1103 // 请先离开已经加入过的房间,才能加入新房间
	RepetitionCreateRoom        = 1104 // 不能重复创建房间
	NotJoinRoom                 = 1105 // 不能加入房间创建房间 稍后再试
	RoomNotExist                = 1106 // 房间不存在
	RoomStatusAbnormal          = 1107 // 房间异常
	JoinRoomErr                 = 1108 // 不能加入已经开始的房间
	HaveAlreadyLeftRoom         = 1109 // 已经离开房间了,不能加入
	RepetitionJoinRoom          = 1110 // 不能重复加入房间
	CrazyLimitAutoBet           = 1111 // boss狂暴阶段 不让自动下注
	NotCrazyStageBet            = 1112 // 不能一键召唤
	GameNotStart                = 1113 // 游戏没开始
	HaveCallSoldiers            = 1114 // 您已经召唤小兵,不能重复召唤
	AlreadySetIsAuto            = 1115 // 当前期 当前层 已经设置过是否自动
	NotAutoUser                 = 1116 // 不是自动用户
	NotCards                    = 1117 // 没有要出的牌
	NotReJoinRoom               = 1118 // 已经不在房间，不能重新加入
	NotRoomOwner                = 1119 // 已经不在房间，不能重新加入
	UpdateCoinErr               = 1120 // 修改积分失败
	NotRoomInviteFriend         = 1121 // 不是房主 不能邀请人
	BeInviteFriendNotLogin      = 1122 // 被邀请人未登陆
	ErrorOperateCardsNum        = 1123 // 出牌数量不对
	IsOutCards                  = 1124 // 已经出过牌
	OutReMakeCardTime           = 1125 // 已超过随牌时间段
	NotExtractCard              = 1126 // 没有可抽的牌
	JoinRoomNotFull             = 1127 // 房间人数不够
	MatchJoinRoomFull           = 1128 // 匹配类型房间只能邀请一个好友
	RoomTypeErr                 = 1129 // 房间类型不匹配
	AuditIng                    = 1130 // 审核中
	AuditPass                   = 1131 // 审核已经通过
	NotAuthRecord               = 1132 // 没有申请纪录
	HaveFriend                  = 1133 // 已经是好友
	AuthHavePass                = 1134 // 已经审核通过
	NotFriendRecord             = 1135 // 不是朋友
	RoomGameStatusIng           = 1136 // 游戏已经开始 不能取消
	NotReadyNotCancel           = 1137 // 没有就绪 不能取消
	NotContinue5                = 1138 // 以抽过卡 不能5连抽
	HaveLikeCard                = 1139 // 已经点过赞
)

// GetErrorMessage 根据错误码 获取错误信息
func GetErrorMessage(code uint32, message string) string {
	var codeMessage string
	codeMap := map[uint32]string{
		WebOK:                      "Success",
		OK:                         "Success",
		NotLoggedIn:                "未登录",
		ParameterIllegal:           "参数不合法",
		UnauthorizedUserID:         "非法的用户ID",
		Unauthorized:               "未授权",
		NotData:                    "没有数据",
		ServerError:                "系统错误",
		ModelAddError:              "添加错误",
		ModelDeleteError:           "删除错误",
		ModelStoreError:            "存储错误",
		OperationFailure:           "操作失败",
		RoutingNotExist:            "路由不存在",
		RepetitiveOperation:        "重复操作",
		SysBusy:                    "系统繁忙 稍后在试",
		Maintenance:                "维护阶段 稍后再来",
		UserScoreNotEnough:         "用户积分不够",
		SelectRoomProhibit:         "该房间已被占用",
		BetLow:                     "押注不正确",
		NoReSelectRoom:             "已经选择过房间",
		ReLogin:                    "已经登陆",
		TokenExpiration:            "身份过期,请重新进入游戏",
		ServerRedisError:           "系统错误 redis",
		GetScoreErr:                "获取远程服务错误",
		DuplicateRequests:          "请求频繁",
		GetRoomConfigErr:           "获取配置房间错误",
		GetCurrTurnErr:             "获取当前轮错误",
		GetTurnDetail:              "获取当前轮房间详情错误",
		ProhibitSelectRoom:         "系统强制匹配时间",
		EnterGameIng:               "繁忙,请重试",
		RepetitionConfirm:          "重复确认",
		DBErr:                      "修改数据失败",
		SCoin:                      "国王积分不足",
		NotSupportCoin:             "不支持的币种",
		ParameterNot:               "参数有误",
		RepetitionFunc:             "重复操作",
		GetDBDataErr:               "获取数据有误",
		UserKnapsackDataEmpty:      "背包数据不足",
		RpcCallRespErr:             "用rpc调用远程服务器错误",
		RpcCallRespDataErr:         "用rpc调用远程服务器返回数据有误",
		ProtocNumberError:          "内部协议号错误",
		DataCompileError:           "数据编译错误",
		GetDataFromDbErr:           "获取DB数据错误",
		AddDataFromDbErr:           "添加数据错误",
		UnauthorizedUserToken:      "非法的用户token",
		HaveBeInviteUser:           "您已经填写过邀请码",
		InviteUserUnreal:           "邀请用户不存在",
		NotInviteYourself:          "不允许邀请自己",
		TavernCreateRoomErr:        "创建房间错误",
		NotLogin:                   "未登陆",
		TavernRoomAlreadyFull:      "已满员",
		JoinRoomFull:               "房间人数已满",
		UserNotInRoom:              "用户不在房间",
		NotCanStartPlay:            "不能开始游戏,请确认用户是否就绪,或者是否达到开始人数",
		BuyRatTimeLimit:            "敬请期待",
		NoWandErr:                  "没有请求需要的标识",
		NoWandPast:                 "需要的标识过期",
		SellOreSinglePriceNotLess0: "单价不能小于0",
		NOtOrder:                   "订单不存在",
		NOtOrderDone:               "订单交易完成",
		ServerTaskIng:              "系统在排队执行中，稍后过来查看",
		NotSelfKickSelf:            "自己不能踢自己",
		NotRoomOwnKick:             "不是房主,不能踢人",
		RoomStatusStopErr:          "房间游戏结束",
		LeavePreRoom:               "请先离开已经加入过的房间",
		RepetitionCreateRoom:       "不能重复创建房间",
		NotJoinRoom:                "不能加入房间创建房间 稍后再试",
		RoomNotExist:               "房间不存在",
		RoomStatusAbnormal:         "房间不存在",
		JoinRoomErr:                "不能加入已经开始的房间",
		HaveAlreadyLeftRoom:        "已经离开房间了,不能加入",
		RepetitionJoinRoom:         "不能重复加入房间",
		GameNotStart:               "游戏没开始",
		NotAutoUser:                "不是自动用户",
		NotCards:                   "没有要出的牌",
		NotReJoinRoom:              "已经不在房间，不能重新加入",
		NotRoomOwner:               "不是房主",
		UpdateCoinErr:              "修改积分失败",
		NotRoomInviteFriend:        "不是房主,不能邀请人",
		BeInviteFriendNotLogin:     "被邀请人未登陆",
		ErrorOperateCardsNum:       "出牌数量不对",
		IsOutCards:                 "已经出过牌",
		OutReMakeCardTime:          "已超过随牌时间段",
		NotExtractCard:             "没有可抽的牌",
		JoinRoomNotFull:            "房间人数不够",
		MatchJoinRoomFull:          "匹配类型房间只能邀请一个好友",
		RoomTypeErr:                "房间类型不匹配",
		AuditIng:                   "审核中",
		AuditPass:                  "审核已经通过",
		NotAuthRecord:              "没有申请纪录",
		HaveFriend:                 "已经是好友",
		AuthHavePass:               "已经审核通过",
		NotFriendRecord:            "不是朋友",
		RoomGameStatusIng:          "游戏已经开始,不能取消",
		NotReadyNotCancel:          "没有就绪,不能取消",
		NotContinue5:               "以抽过卡 不能5连抽",
		HaveLikeCard:               "已经点过赞",
	}

	if message == "" {
		if value, ok := codeMap[code]; ok {
			// 存在
			codeMessage = value
		} else {
			codeMessage = "未定义错误类型!"
		}
	} else {
		codeMessage = message
	}

	return codeMessage
}
