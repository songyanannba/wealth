package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
)

//CREATE TABLE `room_user_num_limit` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`class` int NOT NULL DEFAULT '0' COMMENT '0=人数限制 1=回合数限制',
//`is_show` int unsigned NOT NULL DEFAULT '0' COMMENT '0=显示 1=不显示',
//`num` int NOT NULL DEFAULT '0' COMMENT '房间人数 2/3/4人｜ 回合数 3/5/7',
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='房间配置';

type RoomUserNumLimit struct {
	GVA_MODEL
	Class  int `json:"class" form:"class" gorm:"column:is_show;comment:0=人数限制 1=回合数限制"`
	IsShow int `json:"is_show" form:"is_show" gorm:"column:is_show;comment:0=开放 1=不显示;"`
	Num    int `json:"num" form:"num" gorm:"column:num;comment:房间人数 2/3/4人｜ 回合数 3/5/7;"`
}

// TableName 房间配置
func (o *RoomUserNumLimit) TableName() string {
	return "room_user_num_limit"
}

func GetRoomUserNumLimit() (records []*RoomUserNumLimit, err error) {
	err = global.GVA_SLOT_SERVER_DB.Model(RoomUserNumLimit{}).
		Where("is_show=?", 0).
		Find(&records).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error("sql RoomUserNumLimit error: %s", zap.Error(err))
		return records, err
	}
	return records, nil
}
