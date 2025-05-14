// Package common 通用函数
package common

const (
	WebOK                 = 1    // Success 平台历史原因 原生返回的正确码都是 1
	OK                    = 200  // Success
	NotLoggedIn           = 1000 // 未登录
	ParameterIllegal      = 1001 // 参数不合法
	UnauthorizedUserID    = 1002 // 非法的用户 ID
	Unauthorized          = 1003 // 未授权
	ServerError           = 1004 // 系统错误
	NotData               = 1005 // 没有数据
	ModelAddError         = 1006 // 添加错误
	ModelDeleteError      = 1007 // 删除错误
	ModelStoreError       = 1008 // 存储错误
	OperationFailure      = 1009 // 操作失败
	RoutingNotExist       = 1010 // 路由不存在
	RepetitiveOperation   = 1011 // 重复操作
	SysBusy               = 1012 // 系统繁忙 稍后在试
	Maintenance           = 1013 // 维护阶段 稍后再来
	UserScoreNotEnough    = 1014 // 用户积分不够
	SelectRoomProhibit    = 1015 // 房间被占用
	BetLow                = 1017 // 押注不正确
	NoReSelectRoom        = 1018 // 已经选择过房间
	ReLogin               = 1020
	UserNameRepeat        = 1021 //用户名重复
	TokenExpiration       = 1022
	ServerRedisError      = 1023
	GetScoreErr           = 1024 //获取远程服务错误
	DuplicateRequests     = 1025 //请求频繁
	GetRoomConfigErr      = 1026
	GetCurrTurnErr        = 1027 //获取当前轮错误
	GetTurnDetail         = 1028 //获取当前轮房间详情错误
	ProhibitSelectRoom    = 1029 // 系统强制匹配时间
	EnterGameIng          = 1030 //繁忙,请重试
	NotRegister           = 1031 //未注册
	UnauthorizedUserToken = 1059 // 非法的用户token
	NotLogin              = 1069 // 未登陆
	NotSelfKickSelf       = 1100 // 自己不能踢自己
	NotRoom               = 1032 // 没有房间
	PasswordErr           = 1033 // 密码错误

)

// GetErrorMessage 根据错误码 获取错误信息
func GetErrorMessage(code uint32, message string) string {
	var codeMessage string
	codeMap := map[uint32]string{
		WebOK:                 "Success",
		OK:                    "Success",
		NotLoggedIn:           "未登录",
		ParameterIllegal:      "参数不合法",
		UnauthorizedUserID:    "非法的用户ID",
		Unauthorized:          "未授权",
		NotData:               "没有数据",
		ServerError:           "系统错误",
		ModelAddError:         "添加错误",
		ModelDeleteError:      "删除错误",
		ModelStoreError:       "存储错误",
		OperationFailure:      "操作失败",
		RoutingNotExist:       "路由不存在",
		RepetitiveOperation:   "重复操作",
		SysBusy:               "系统繁忙,稍后在试",
		Maintenance:           "维护阶段,稍后再来",
		UserScoreNotEnough:    "用户积分不够",
		SelectRoomProhibit:    "该房间已被占用",
		BetLow:                "押注不正确",
		NoReSelectRoom:        "已经选择过房间",
		ReLogin:               "已经登陆",
		TokenExpiration:       "身份过期,请重新进入游戏",
		ServerRedisError:      "系统错误 redis",
		GetScoreErr:           "获取远程服务错误",
		DuplicateRequests:     "请求频繁",
		GetRoomConfigErr:      "获取配置房间错误",
		GetCurrTurnErr:        "获取当前轮错误",
		GetTurnDetail:         "获取当前轮房间详情错误",
		ProhibitSelectRoom:    "系统强制匹配时间",
		EnterGameIng:          "繁忙,请重试",
		UserNameRepeat:        "name repeat",
		NotRegister:           "not register",
		UnauthorizedUserToken: "非法的用户token",
		NotLogin:              "not login",
		NotRoom:               "not room",
		PasswordErr:           "密码错误",
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
