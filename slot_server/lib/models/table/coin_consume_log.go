package table

import (
	"go.uber.org/zap"
	"slot_server/lib/global"

	"time"
)

//CREATE TABLE `coin_consume_log` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '会员ID',
//`op_type` int NOT NULL DEFAULT '0' COMMENT '0:默认 1=开包 2=重随',
//`curr_num` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '币数量',
//`before_num` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '前-货币',
//`after_num` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '后-货币',
//`desc` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '描述',
//`date` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '年月日',
//`date_time` datetime DEFAULT NULL,
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT NULL,
//PRIMARY KEY (`id`),
//KEY `user_id` (`user_id`) USING BTREE
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='金币扣减纪录';

type CoinConsumeLog struct {
	GVA_MODEL
	UserId    string     `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id;"`
	OpType    int        `json:"op_type" form:"op_type" gorm:"column:op_type;comment:0:默认 1=开包 2=重随;"`
	CurrNum   float64    `json:"curr_num" form:"curr_num" gorm:"column:curr_num;comment:币数量;"`
	AfterNum  float64    `json:"after_num" form:"after_num" gorm:"column:after_num;comment:后-货币;"`
	BeforeNum float64    `json:"before_num" form:"before_num" gorm:"column:before_num;comment:前-货币;"`
	Desc      string     `json:"desc" form:"desc" gorm:"column:desc;comment:描述;"`
	Date      string     `json:"date" form:"date" gorm:"column:date;comment:年月日;"`
	DateTime  *time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
}

// TableName 挖矿每日详情 日志log
func (o *CoinConsumeLog) TableName() string {
	return "coin_consume_log"
}

func CreateCoinConsumeLog(record *CoinConsumeLog) error {
	err := global.GVA_SLOT_SERVER_DB.Model(CoinConsumeLog{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql CoinConsumeLog error: %s", zap.Error(err))
		return err
	}
	return nil
}
