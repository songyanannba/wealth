package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"time"
)

//CREATE TABLE `meme_room` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '房间创建者',
//`owner` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '房主',
//`room_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '房间号',
//`name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '名称',
//`desc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '描述',
//`is_open` tinyint(1) NOT NULL DEFAULT '0' COMMENT '房间状态: 1=开放中,2=已满员,3=已解散,4=进行中,5=已结束 6=异常房间 7=服务字段清理残存房间 8=清理匹配成功用户之前的房间 ',
//`date` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '年月日',
//`user_num_limit` int NOT NULL DEFAULT '0' COMMENT '用户人数限制 2人场 3 人场 4人场',
//`room_turn_num` int NOT NULL DEFAULT '0' COMMENT '房间回合数 3/5/7，默认5',
//`room_type` tinyint NOT NULL DEFAULT '0' COMMENT '房间类型 1=好友约战 2=匹配',
//`room_class` tinyint NOT NULL DEFAULT '0' COMMENT '0:快速匹配房间 1:正常匹配房间（暂时没用到）',
//`room_level` tinyint NOT NULL DEFAULT '0' COMMENT '房间 等级 0 初级 1 中级 2 高级（暂时没用到）',
//`is_go_on` int NOT NULL DEFAULT '0' COMMENT '0 默认 ；1=继续游戏的房间',
//`date_time` datetime DEFAULT NULL,
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`),
//KEY `room_no` (`room_no`) USING BTREE,
//KEY `owner` (`owner`) USING BTREE,
//KEY `date` (`date`) USING BTREE
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='meme房间';

const TavernRoomOpen = 1     //开放中
const TavernRoomFullOpen = 2 //已满员
const TavernRoomNoOpen = 3   //已关闭

const RoomTypeCard = 1 //好友约战

// RoomLevel 房间等级
type RoomLevel int8

const (
	RoomLevelPri RoomLevel = 0
	RoomLevelMid RoomLevel = 1
	RoomLevelExp RoomLevel = 2
)

const (
	RoomStatusOpen         = 1 //开放中 未满员
	RoomStatusFill         = 2 //满员
	RoomStatusDissolve     = 3 //已解散
	RoomStatusIng          = 4 //进行中
	RoomStatusStop         = 5 //已结束
	RoomStatusAbnormal     = 6 //异常房间
	RoomClearRoom          = 7 //清理房间
	RoomClearMatchSuccUser = 8 //清理匹配成功用户之前的房间
)

func IsInRoomLevel(roomLevel RoomLevel) bool {
	if roomLevel != RoomLevelPri && roomLevel != RoomLevelMid && roomLevel != RoomLevelExp {
		return false
	}
	return true
}

type RoomClass int8

const (
	RoomClassInvite = 0
	RoomClassMatch  = 1 //继续游戏
)

type RoomType int8

const (
	RoomTypeYueZhan = 1 //好友约战
	RoomTypeMatch   = 2 //匹配模式
)

type AnimalPartyRoom struct {
	GVA_MODEL
	UserId       string     `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id;"`
	Owner        string     `json:"owner" form:"owner" gorm:"column:owner;comment:房主;"`
	RoomNo       string     `json:"room_no" form:"room_no" gorm:"column:room_no;comment:房间号;"`
	Name         string     `json:"name" form:"name" gorm:"column:name;comment:名称;"`
	Desc         string     `json:"desc" form:"desc" gorm:"column:desc;comment:描述;"`
	IsOpen       int8       `json:"is_open" form:"is_open" gorm:"column:is_open;default:0;comment:房间状态:1=开放中,2=已满员,3=已解散,4=进行中,5=已结束,6=异常房间 8=清理匹配成功用户之前的房间;"`
	RoomType     int8       `json:"room_type" form:"room_type" gorm:"column:room_type;default:0;comment:房间类型 1=动物派对（全局）;"`
	RoomLevel    int8       `json:"room_level" form:"room_level" gorm:"column:room_level;default:0;comment:房间等级 0 初级 1 中级 2 高级;"`
	RoomClass    int8       `json:"room_class" form:"room_class" gorm:"column:room_class;default:0;comment:0: 1:继续游戏;"`
	UserNumLimit int        `json:"user_num_limit" form:"user_num_limit" gorm:"column:user_num_limit;default:0;comment:用户人数限制 2人场 3 人场 4人场;"`
	RoomTurnNum  int        `json:"room_turn_num" form:"room_turn_num" gorm:"column:room_turn_num;default:0;comment:房间回合数 3/5/7，默认5"`
	IsGoOn       int        `json:"is_go_on" form:"is_go_on" gorm:"column:is_go_on;default:0;comment:0 默认；1= 继续游戏的房间;"`
	CityId       int        `json:"city_id" form:"city_id" gorm:"column:city_id;default:0;comment:0 默认;"`
	Date         string     `json:"date" form:"date" gorm:"column:date;comment:年月日;"`
	Period       string     `json:"period" form:"period" gorm:"column:period;default:0;comment:当前第几期"`
	DateTime     *time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
}

// TableName 用户房间
func (o *AnimalPartyRoom) TableName() string {
	return "animal_party_room"
}

func NewAnimalPartyRoom(userID, owner, roomNo, name, desc, period string, isOpen, roomType, roomLevel, roomClass int8, turnNum, userNumLimit int) *AnimalPartyRoom {
	return &AnimalPartyRoom{
		UserId:       userID,
		Owner:        owner,
		RoomNo:       roomNo,
		Name:         name,
		Desc:         desc,
		IsOpen:       isOpen,
		RoomType:     roomType,
		RoomLevel:    roomLevel,
		RoomClass:    roomClass,
		RoomTurnNum:  turnNum, //回合数
		UserNumLimit: userNumLimit,
		Date:         helper.YearMonthDayStr(),
		DateTime:     helper.LocalTime(),
	}
}

func CreateMemeRoom(record *AnimalPartyRoom) error {
	err := global.GVA_SLOT_SERVER_DB.Model(AnimalPartyRoom{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql CreateMemeRoom error: %s", zap.Error(err))
		return err
	}
	return nil
}

// DelMemeRoom 删除房间
func DelMemeRoom(roomNo string) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(AnimalPartyRoom{}).
		Where("room_no = ?", roomNo).
		Delete(AnimalPartyRoom{}).
		Error
	if err != nil {
		global.GVA_LOG.Error("sql DelMemeRoom error: %s", zap.Error(err))
		return err
	}
	return nil
}

func SaveMemeRoom(record *AnimalPartyRoom) error {
	err := global.GVA_SLOT_SERVER_DB.Model(AnimalPartyRoom{}).
		Where("id = ?", record.ID).
		Save(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql MemeRoom error: %s", zap.Error(err))
		return err
	}
	return nil
}

func GetMemeRoomById(id int) (record *AnimalPartyRoom, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(AnimalPartyRoom{}).
		Where("id = ?", id).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return
	}
	return record, nil
}

func GetMemeRoomByIdDesc() (record *AnimalPartyRoom, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(AnimalPartyRoom{}).
		Order("id desc").
		First(&record).
		Limit(1).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return
	}
	return record, nil
}

func SlotRoomByRoomNo(roomNo string) (record *AnimalPartyRoom, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(AnimalPartyRoom{}).
		Where("room_no = ?", roomNo).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

// GetMemeRoomByUid 返回用户没有结束的房间
func GetMemeRoomByUid(userID string) (record *AnimalPartyRoom, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(AnimalPartyRoom{}).
		Where("user_id = ? and is_open in ? ", userID, []int{RoomStatusOpen, RoomStatusFill, RoomStatusIng}).
		Order("id desc").
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func GetMemeRoomByOwner(userID string) (record *AnimalPartyRoom, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(AnimalPartyRoom{}).
		Where("owner = ? and is_open in ? ", userID, []int{RoomStatusOpen, RoomStatusFill, RoomStatusIng}).
		Order("id desc").
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func GetMemeRoomPage(roomId int, roomLevel RoomLevel) (val []*AnimalPartyRoom, err error) {
	db := global.GVA_SLOT_SERVER_DB.Model(AnimalPartyRoom{})
	if roomId == 0 {
		if roomLevel == RoomLevelPri || roomLevel == RoomLevelMid || roomLevel == RoomLevelExp {
			db.Where("room_level = ? ", roomLevel).Where("is_open = ?", RoomStatusOpen).Where("room_class = ?", 1).Where("is_go_on = ?", 0)
		} else {
			db.Where("is_open = ?", RoomStatusOpen).Where("room_class = ?", 1).Where("is_go_on = ?", 0)
		}
	} else {
		if roomLevel == RoomLevelPri || roomLevel == RoomLevelMid || roomLevel == RoomLevelExp {
			db.Where("room_level = ? ", roomLevel).Where("is_open = ? and id > ?", RoomStatusOpen, roomId).Where("room_class = ?", 1).Where("is_go_on = ?", 0)
		} else {
			db.Where("is_open = ? and id > ?", RoomStatusOpen, roomId).Where("room_class = ?", 1).Where("is_go_on = ?", 0)
		}
	}

	err = db.Limit(10).Find(&val).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return nil, err
	}
	return val, nil
}

func MemeRoomPageIsNext(roomId int) (count int64, err error) {
	db := global.GVA_SLOT_SERVER_DB.
		Model(AnimalPartyRoom{})

	if roomId != 0 {
		db.Where(" id > ?", roomId)
	}

	err = db.Count(&count).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return count, err
	}
	return count, nil
}

func NewestMemeRoomByUid(userId string) (record *AnimalPartyRoom, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(AnimalPartyRoom{}).
		Where("user_id = ?", userId).
		Order("id desc").
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func NewestMemeRoomByRoomNo(roomNo string) (record *AnimalPartyRoom, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(AnimalPartyRoom{}).
		Where("room_no = ?", roomNo).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func NewestNormalMemeRoomByRoomNo(roomNo string) (record *AnimalPartyRoom, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(AnimalPartyRoom{}).
		Where("room_no = ? and is_open in ?", roomNo, []int{RoomStatusOpen, RoomStatusFill, RoomStatusIng}).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}
