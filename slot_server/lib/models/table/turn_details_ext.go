package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"time"
)

//CREATE TABLE `turn_details_ext` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`details_id` int NOT NULL DEFAULT '0' COMMENT '详情ID',
//`room_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间号',
//`turn` int NOT NULL DEFAULT '0' COMMENT '轮',
//`issue_desc` varchar(1000) DEFAULT NULL COMMENT '问题描述',
//`cards` json DEFAULT NULL COMMENT '剩余牌',
//`out_cards` json DEFAULT NULL COMMENT '出牌',
//`date` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '年月日',
//`date_time` datetime DEFAULT NULL,
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`),
//KEY `details_id` (`details_id`) USING BTREE,
//KEY `room_no` (`room_no`) USING BTREE
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='游戏局数结算详情扩展';

type TurnDetailsExt struct {
	GVA_MODEL
	DetailsId int        `json:"details_id" form:"details_id" gorm:"column:details_id;default:0;comment:详情ID;"`
	RoomNo    string     `json:"room_no" form:"room_no" gorm:"column:room_no;comment:房间号;"`
	IssueDesc string     `json:"issue_desc" form:"issue_desc" gorm:"column:issue_desc;comment:问题描述;"`
	Cards     string     `json:"cards" form:"cards" gorm:"column:cards;comment:剩余牌;"`
	OutCards  string     `json:"out_cards" form:"out_cards" gorm:"column:out_cards;comment:出牌;"`
	Date      string     `json:"date" form:"date" gorm:"column:date;comment:年月日;"`
	DateTime  *time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
	Turn      int        `json:"turn" form:"turn" gorm:"column:turn;default:0;comment:轮;"`
}

// TableName 游戏局数结算详情
func (o *TurnDetailsExt) TableName() string {
	return "turn_details_ext"
}

func NewTurnDetailsExt(detailsId int, roomNo string, turn int, cards, outCards, IssueDesc string) TurnDetailsExt {
	return TurnDetailsExt{
		RoomNo:    roomNo,
		DetailsId: detailsId,
		Turn:      turn,
		Cards:     cards,
		OutCards:  outCards,
		IssueDesc: IssueDesc,
		Date:      helper.YearMonthDayStr(),
		DateTime:  helper.LocalTime(),
	}
}

func CreateTurnDetailsExt(record *TurnDetailsExt) error {
	err := global.GVA_SLOT_SERVER_DB.Model(TurnDetailsExt{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql CreateTurnDetailsExt error: %s", zap.Error(err))
		return err
	}
	return nil
}

func GetTurnDetailsExtByDetailsId(detailsId int) (record *TurnDetailsExt, err error) {
	err = global.GVA_SLOT_SERVER_DB.Model(TurnDetailsExt{}).
		Where("details_id = ? ", detailsId).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(" sql GetTurnDetailsExtByDetailsId error: %s", zap.Error(err))
		return record, err
	}
	return record, nil
}

func GetTurnDetailsExt(roomNo string) (records []*TurnDetailsExt, err error) {
	err = global.GVA_SLOT_SERVER_DB.Model(TurnDetailsExt{}).
		Where("room_no = ?", roomNo).
		Find(&records).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(" GetTurnDetailsExt error: %s", zap.Error(err))
		return records, err
	}
	return records, nil
}

func SaveTurnDetailsExt(record *TurnDetailsExt) error {
	err := global.GVA_SLOT_SERVER_DB.Model(TurnDetailsExt{}).
		Where("id = ?", record.ID).
		Save(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("sql SaveTurnDetailsExt error: %s", zap.Error(err))
		return err
	}
	return nil
}

func UpdateTurnDetailsExt(detailsId int, id int, values map[string]interface{}) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(TurnDetailsExt{}).
		Where("details_id = ? and id = ?", detailsId, id).
		Updates(values).
		Error
	if err != nil {
		global.GVA_LOG.Error("UpdateTurnDetailsExt error: %s", zap.Error(err))
		return err
	}
	return nil
}
