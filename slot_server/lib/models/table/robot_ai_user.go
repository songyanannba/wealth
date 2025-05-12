package table

import (
	"gorm.io/gorm"
	"slot_server/lib/global"
)

//CREATE TABLE `robot_ai_user` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`class` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1:机器人1 2:机器人2 3:机器人3',
//`stage` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1:随牌阶段 2:出牌阶段 3:点赞阶段',
//`random_num` varchar(50) NOT NULL DEFAULT '' COMMENT '随机次数',
//`out_time` varchar(50) NOT NULL DEFAULT '' COMMENT '出牌时间',
//`out_theme` tinyint(1) NOT NULL DEFAULT '0' COMMENT '出牌策略 1:随机 2:等级最高 ',
//`like_time` varchar(50) NOT NULL DEFAULT '' COMMENT '点赞时间',
//`like_theme` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1:随机  2:牌等级优先 3:跟风',
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='机器人配置';

type RobotAiUser struct {
	GVA_MODEL
	Class     int8   `json:"class" form:"class" gorm:"column:class;default:0;comment:1:机器人1 2:机器人2 3:机器人3;"`
	Stage     int8   `json:"stage" form:"stage" gorm:"column:stage;default:0;comment:1:机器人1 2:机器人2 3:机器人3;"`
	RandomNum string `json:"random_num" form:"random_num" gorm:"column:random_num;default:0;comment:随机次数;"`
	OutTime   string `json:"out_time" form:"out_time" gorm:"column:out_time;default:0;comment:出牌时间;"`
	OutTheme  int8   `json:"out_theme" form:"out_theme" gorm:"column:out_theme;default:0;comment:出牌策略 1:随机 2:等级最高;"`
	LikeTime  string `json:"like_time" form:"like_time" gorm:"column:like_time;default:0;comment:点赞时间;"`
	LikeTheme int8   `json:"like_theme" form:"like_theme" gorm:"column:like_theme;default:0;comment:1:随机 2:牌等级优先 3:跟风;"`
}

// TableName 用户房间
func (o *RobotAiUser) TableName() string {
	return "robot_ai_user"
}

// GetRobotAiUser 获取机器人
func GetRobotAiUser() (record []*RobotAiUser, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(RobotAiUser{}).
		Find(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}
