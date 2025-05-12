package logic

import (
	"encoding/json"
	"go.uber.org/zap"
	"slot_server/lib/config"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models"
	"slot_server/lib/models/table"
	"slot_server/lib/utils/cache"
)

func GetMbCardConfig(version int) []*table.MbCardConfig {
	cardConfigs := make([]*table.MbCardConfig, 0)
	cardConfigStr, err := cache.GetMbCardConfig(helper.Itoa(version))
	if err != nil {
		global.GVA_LOG.Error(" GetMbCardConfig fail", zap.Error(err))
		return cardConfigs
	}

	if len(cardConfigStr) == 0 {
		//db
		cardConfigs, err = table.GetMbCardConfigByVersion(version)
		if err != nil {
			global.GVA_LOG.Error(" GetMbCardConfigByVersion", zap.Error(err))
			return cardConfigs
		}

		cardConfigMarshal, _ := json.Marshal(cardConfigs)
		err = cache.SetMbCardConfig(helper.Itoa(version), string(cardConfigMarshal))
		if err != nil {
			global.GVA_LOG.Error(" GetMbCardConfig SetMbCardConfig", zap.Error(err))
		}
	} else {
		err = json.Unmarshal([]byte(cardConfigStr), &cardConfigs)
		if err != nil {
			global.GVA_LOG.Error(" GetMbCardConfig Unmarshal fail", zap.Error(err))
			return cardConfigs
		}
	}
	return cardConfigs
}

// UnpackCardVersionAndNum 拆包 + 缓存
func UnpackCardVersionAndNum(userId string, version, num int) [][]*models.HandListCard {
	unpackCardNum := make([][]*models.HandListCard, 0)
	cardConfigs := make([]*table.MbCardConfig, 0)

	cardConfigs = GetMbCardConfig(version)
	if len(cardConfigs) == 0 {
		return unpackCardNum
	}

	if num == config.UnpackCardNum1 {
		unpackCard := make([]*models.HandListCard, 0)
		if len(cardConfigs) > config.UnpackCardCount {
			helper.SliceShuffle(cardConfigs)
		}

		for k, confItem := range cardConfigs {
			if k > config.UnpackCardCount {
				break
			}
			unpackCard = append(unpackCard, &models.HandListCard{
				CardId: confItem.ID,
				Name:   confItem.Name,
				Suffix: confItem.SuffixName,
				Level:  confItem.Level,
			})
		}
		unpackCardNum = append(unpackCardNum, unpackCard)
	}

	//暂时没有开10包的需求
	//if num == 5 || num == config.UnpackCardNum10 {

	//首次赠送
	if num == 5 {
		for i := 0; i < num; i++ {
			unpackCard := make([]*models.HandListCard, 0)
			if len(cardConfigs) > config.UnpackCardCount {
				helper.SliceShuffle(cardConfigs)
			}

			for k, confItem := range cardConfigs {
				if k > config.UnpackCardCount {
					break
				}
				unpackCard = append(unpackCard, &models.HandListCard{
					CardId: confItem.ID,
					Name:   confItem.Name,
					Suffix: confItem.SuffixName,
					Level:  confItem.Level,
				})
			}
			unpackCardNum = append(unpackCardNum, unpackCard)
		}
	}

	//开包以后放入图鉴
	for _, unpackCardItem := range unpackCardNum {
		BatchUpdateUserHandbook(userId, unpackCardItem)
	}

	return unpackCardNum
}

func CardConfList(userId string, lastId, level int) []*models.HandListCard {
	cardConfList := make([]*models.HandListCard, 0)

	cardConfigStr, err := cache.GetMbCardConfigPage(lastId, level)
	if err != nil {
		global.GVA_LOG.Error("CardConfList fail", zap.Error(err))
		return cardConfList
	}

	configs := make([]*table.MbCardConfig, 0)
	if len(cardConfigStr) == 0 {
		//db
		configs, err = table.GetMbCardConfigByLastId(lastId, level)
		if err != nil {
			global.GVA_LOG.Error("CardConfList GetMbCardConfigByLastId", zap.Error(err))
			return cardConfList
		}

		if len(configs) > 0 {
			cardConfigMarshal, _ := json.Marshal(configs)
			err = cache.SetMbCardConfigPage(lastId, level, string(cardConfigMarshal))
			if err != nil {
				global.GVA_LOG.Error("UnpackCardVersionAndNum SetMbCardConfig", zap.Error(err))
			}
		}
	} else {
		err = json.Unmarshal([]byte(cardConfigStr), &configs)
		if err != nil {
			global.GVA_LOG.Error("UnpackCardVersionAndNum Unmarshal fail", zap.Error(err))
			return cardConfList
		}
	}

	for _, confItem := range configs {
		isOwn := false

		recordVal, err := cache.GetMbUserHandbook(userId, confItem.ID)
		if err != nil {
			global.GVA_LOG.Error("CardConfList GetMbUserHandbook", zap.Error(err))
			continue
		}

		userHandbook := &table.MbUserHandbook{}
		if len(recordVal) > 0 {
			err = json.Unmarshal([]byte(recordVal), &userHandbook)
			if err != nil {
				global.GVA_LOG.Error("UnpackCardVersionAndNum Unmarshal fail", zap.Error(err))
				return cardConfList
			}
		} else {
			userHandbook, err = table.GetUserHandbookByCardId(userId, confItem.ID)
			if err != nil {
				global.GVA_LOG.Error("CardConfList GetUserHandbookByCardId", zap.Error(err))
				continue
			}

			if userHandbook.ID > 0 {
				recordMarshal, _ := json.Marshal(userHandbook)
				err = cache.SetMbUserHandbook(userId, confItem.ID, string(recordMarshal))
				if err != nil {
					global.GVA_LOG.Error("CardConfList SetMbUserHandbook", zap.Error(err))
				}
			}
		}

		if userHandbook.ID > 0 {
			isOwn = true
		}

		cardConfList = append(cardConfList, &models.HandListCard{
			CardId: confItem.ID,
			Name:   confItem.Name,
			Suffix: confItem.SuffixName,
			IsOwn:  isOwn,
			Level:  confItem.Level,
		})
	}

	return cardConfList
}

func CardConfListAndIsHaveNextPage(userId string, lastId, level int) ([]*models.HandListCard, bool, int) {
	pageNum := 20
	isHaveNextPage := false
	pageCount, err := table.CardConfigPageIsNext(lastId, level)
	if err != nil {
		global.GVA_LOG.Error("CardConfListAndIsHaveNextPage CardConfigPageIsNext :", zap.Error(err))
	}
	if pageCount > int64(pageNum) {
		isHaveNextPage = true
	}

	//对应类型一共多少
	allCount, err := table.CardConfigCount(level)
	if err != nil {
		global.GVA_LOG.Error("CardConfListAndIsHaveNextPage CardConfigCount :", zap.Error(err))
	}

	return CardConfList(userId, lastId, level), isHaveNextPage, int(allCount)
}

// UnpackCardVersion 开包 根据对应版本
//func UnpackCardVersion(userId string, version, num int) []*models.HandListCard {
//	unpackCard := make([]*models.HandListCard, 0)
//
//	//todo 加缓存
//	configs, err := table.GetMbCardConfigByVersion(version)
//	if err != nil {
//		global.GVA_LOG.Error("UnpackCardVersion GetMbCardConfigByVersion", zap.Error(err))
//		return unpackCard
//	}
//
//	if len(configs) > config.UnpackCardCount {
//		helper.SliceShuffle(configs)
//	}
//
//	for k, confItem := range configs {
//		if k > config.UnpackCardCount {
//			break
//		}
//		unpackCard = append(unpackCard, &models.HandListCard{
//			CardId: confItem.ID,
//			Name:   confItem.Name,
//			Suffix: confItem.SuffixName,
//		})
//	}
//
//	//开包以后放入图鉴
//	BatchUpdateUserHandbook(userId, unpackCard)
//	return unpackCard
//}
