package common

type OperateCardType int

// 0:看牌 1:出牌 2:表情 3:重随
const (
	LookCards OperateCardType = iota
	OperateCards
	OpeEmoji
	ReMakeCards
)
