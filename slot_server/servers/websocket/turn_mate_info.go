package websocket

import "slot_server/lib/models"

// SetLikeCardsCard 本轮点赞牌的集合
func (rs *ComRoomSpace) SetLikeCardsCard(LikeCards *models.LikeCard) {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()

	turn := rs.GetTurn()

	//收集本轮所有被点赞的牌
	cards, ok := rs.TurnMateInfo.likeCards[turn]
	if !ok {
		rs.TurnMateInfo.likeCards[turn] = make([]*models.LikeCard, 0)
	}

	isBeLiked := false
	for _, card := range cards {
		//同一个用户 相同的牌被多次点赞
		if card.CardId == LikeCards.CardId && LikeCards.LikeUserId == card.LikeUserId {
			isBeLiked = true
			card.LikeNum += 1
		}
	}

	if !isBeLiked {
		LikeCards.LikeNum += 1
		rs.TurnMateInfo.likeCards[turn] = append(rs.TurnMateInfo.likeCards[turn], LikeCards)
	}
}

// GetCurrTurnLikeCards 获取本轮点赞集合
func (rs *ComRoomSpace) GetCurrTurnLikeCards() []*models.LikeCard {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()

	res := make([]*models.LikeCard, 0)
	turn := rs.GetTurn()

	cards, ok := rs.TurnMateInfo.likeCards[turn]
	if !ok {
		return res
	}
	return cards
}

func (rs *ComRoomSpace) TurnLikeCards(turn int) []*models.LikeCard {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()

	res := make([]*models.LikeCard, 0)
	cards, ok := rs.TurnMateInfo.likeCards[turn]
	if !ok {
		return res
	}
	return cards
}

// AllTurnLikeCards 获取全部轮的点赞集合
func (rs *ComRoomSpace) AllTurnLikeCards() []*models.LikeCard {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()

	res := make([]*models.LikeCard, 0)

	for _, allCards := range rs.TurnMateInfo.likeCards {
		res = append(res, allCards...)
	}
	return res
}

// IsAllUserLikeCard 判断是否都点赞
func (rs *ComRoomSpace) IsAllUserLikeCard() bool {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()
	turn := rs.GetTurn()

	cards, ok := rs.TurnMateInfo.likeCards[turn]
	if !ok {
		return false
	}

	likeNum := 0
	for _, card := range cards {
		likeNum += card.LikeNum
	}

	//点在数和用户一致 说明所有的用户已经点赞
	if likeNum != len(rs.UserInfos) {
		return false
	}

	return true
}
