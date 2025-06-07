package models

type UserInfoResp struct {
	UserName string  `json:"user_name"` // 用户名称
	Amount   float64 `json:"amount"`    // 金额
	City     string  `json:"city"`      // 城市
}
