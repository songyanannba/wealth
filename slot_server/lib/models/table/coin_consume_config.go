package table

import (
	"gorm.io/gorm"
	"slot_server/lib/global"
)

//CREATE TABLE `coin_consume_config` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`type` int NOT NULL DEFAULT '0' COMMENT '分类 1=开包 2=重随',
//`coin_num` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '金币消耗量',
//`desc` varchar(300) DEFAULT NULL COMMENT '介绍',
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='积分消耗配置';

type CoinConsume int

var (
	CoinConsumeUnpackCard CoinConsume = 1 //开包消耗
	CoinConsumeReCard1    CoinConsume = 2 //重随第1次消耗
	CoinConsumeReCard2    CoinConsume = 3 //重随第2次消耗
	CoinConsumeReCard3    CoinConsume = 4 //重随第大于等于3次消耗
)

type CoinConsumeConfig struct {
	GVA_MODEL
	Type    int     `json:"type" form:"type" gorm:"column:type;default:0;comment:分类 1=开包 2=重随"`
	CoinNum float64 `json:"coin_num" form:"coin_num" gorm:"column:coin_num;comment:金币消耗量;"`
	Desc    string  `json:"desc" form:"date" gorm:"column:desc;comment:介绍;"`
}

func (o *CoinConsumeConfig) TableName() string {
	return "coin_consume_config"
}

func GetCoinConsume(num int) CoinConsume {
	if num == 1 {
		return CoinConsumeReCard1
	} else if num == 2 {
		return CoinConsumeReCard2
	} else {
		return CoinConsumeReCard3
	}
}

// GetCoinConsumeConfigByType 获取消耗
func GetCoinConsumeConfigByType(typ int) (record *CoinConsumeConfig, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(CoinConsumeConfig{}).
		Where("type = ?", typ).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}
