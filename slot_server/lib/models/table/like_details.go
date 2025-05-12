package table

import "time"

//CREATE TABLE `like_details` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`turn` int NOT NULL DEFAULT '0' COMMENT '游戏局数',
//`user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户ID',
//`nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称',
//`room_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间号',
//`like_num` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0',
//`cards` json DEFAULT NULL COMMENT '被点赞的牌',
//`date` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '年月日',
//`date_time` datetime DEFAULT NULL,
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`),
//KEY `room_no` (`room_no`) USING BTREE,
//KEY `user_id` (`user_id`) USING BTREE,
//KEY `date` (`date`) USING BTREE
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='点赞详情';

type LikeDetails struct {
	GVA_MODEL
	Turn     int        `json:"turn" form:"turn" gorm:"column:turn;default:0;comment:第几轮"`
	UserId   string     `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id;"`
	RoomNo   string     `json:"room_no" form:"room_no" gorm:"column:room_no;comment:房间号;"`
	Nickname string     `json:"nickname" form:"nickname" gorm:"column:nickname;comment:昵称;"`
	LikeNum  int        `json:"like_num" form:"like_num" gorm:"column:like_num;default:0;comment:点赞数量"`
	Cards    string     `json:"cards" form:"cards" gorm:"column:cards;comment:被点赞的牌;"`
	Date     string     `json:"date" form:"date" gorm:"column:date;comment:年月日;"`
	DateTime *time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
}

// TableName 游戏玩家表
func (o *LikeDetails) TableName() string {
	return "like_details"
}
