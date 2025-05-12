package websocket

import (
	"github.com/google/uuid"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models"
	"time"
)

func (trMgr *roomManager) DealMatchIngUserAddRobot() {
	//超过20秒 还没匹配成功 加入机器人逻辑
	//单排加3个机器人 双排加2个机器人
	match4 := 4
	matchIngRoom1User := trMgr.MatchIngRoom.MatchIngRoom1User

	//当前匹配中的单排人数（不够4人）
	matchIngRoom1UserLen := len(matchIngRoom1User)

	if matchIngRoom1UserLen <= 0 || matchIngRoom1UserLen > match4 {
		return
	}

	//当前需要加机器人的用户
	matchIngRoom1UserGroup := make([]*MatchIngRoomInfo, 0)
	matchIngRoom1UserGroup = matchIngRoom1User[0:matchIngRoom1UserLen]
	trMgr.MatchIngRoom.MatchIngRoom1User = matchIngRoom1User[matchIngRoom1UserLen:]

	//前两个用户匹配成功
	matchUser := &MatchGroupRoomInfo{}
	//匹配成功 多少对
	matchRoomUserArr := make([]*MatchGroupRoomInfo, 0)

	global.GVA_LOG.Infof("DealMatchIngUserAddRobot 单排加机器人 {%v}人场匹配期...", match4)

	for _, matchRoomInfo := range matchIngRoom1UserGroup {
		if len(matchRoomInfo.UserInfoArr) != 1 {
			//没有足够的匹配用户
			global.GVA_LOG.Infof("DealMatchIngUserAddRobot RoomNo:%v, UserInfoMapLen:%v", matchRoomInfo.RoomNo, len(matchRoomInfo.UserInfoArr))
			continue
		}

		//匹配时间 在20秒内 不加机器人
		if helper.LocalTime().Before(matchRoomInfo.StartTime.Add(20 * time.Second)) {
			trMgr.MatchIngRoom.MatchIngRoom1User = append(trMgr.MatchIngRoom.MatchIngRoom1User, matchRoomInfo)
			continue
		}

		matchUser.DelRoomNo = append(matchUser.DelRoomNo, matchRoomInfo.RoomNo)

		//匹配用户 顺便过滤已经退出的用户
		for k, _ := range matchRoomInfo.UserInfoArr {
			uInfo := matchRoomInfo.UserInfoArr[k]
			if len(matchUser.UserInfoArr) < match4 {
				matchUser.UserInfoArr = append(matchUser.UserInfoArr, uInfo)
			}

			//加三个机器人
			robotUser := trMgr.GetRobot(3)
			matchUser.UserInfoArr = append(matchUser.UserInfoArr, robotUser...)

			if len(matchUser.UserInfoArr) == match4 {
				newMatchUser := MatchGroupRoomInfo{
					RoomNo:      "",
					DelRoomNo:   matchUser.DelRoomNo,
					UserInfoArr: matchUser.UserInfoArr,
				}
				matchRoomUserArr = append(matchRoomUserArr, &newMatchUser)
				//重置
				matchUser = &MatchGroupRoomInfo{}
			}
		}
	}

	trMgr.MatchGroupStart(matchRoomUserArr) //1+3机器人
	global.GVA_LOG.Infof("DealMatchIngUserAddRobot 单排加机器人 {%v}人场匹配 结束...", match4)
}

func (trMgr *roomManager) GetRobot(num int) []*models.UserInfo {
	userInfoArr := make([]*models.UserInfo, 0)
	//创建n 个机器人

	for i := 0; i < num; i++ {
		userID := uuid.New().String()
		userInfo := &models.UserInfo{
			UserID:   userID,
			Nickname: GetRobotNickname(),
			UserProperty: models.UserProperty{
				IsRobot:    1,
				RobotClass: RobotClass(),
			},
			UserExt: models.UserExt{},
		}
		userInfoArr = append(userInfoArr, userInfo)
	}
	return userInfoArr
}

func GetRobotNickname() string {
	var name string
	nameArr := []string{
		"退网局常驻",
		"老八秘制小卡",
		"我出布你出寄",
		"卡比巴拉海",
		"有狗在偷窥",
		"在逃公主",
		"卡牌刺客",
		"赛博菩萨显灵了",
		"卡姿兰大法师",
		"鼠鼠我鸭",
		"薯条杀手",
		"啊？尊嘟假嘟",
		"网瘾诱捕器",
		"电子咸鱼腌入味",
		"我佛糍粑",
		"九转大肠仙人",
		"雪豹闭嘴",
		"泰裤辣条",
		"急急国王的披风",
		"芝士雪豹",
	}
	name = nameArr[helper.RandInt(len(nameArr))]
	if name == "" {
		name = nameArr[0]
	}
	return name
}

func RobotClass() int {
	var class int
	nameArr := []int{1, 2, 3}
	class = nameArr[helper.RandInt(len(nameArr))]
	if class == 0 {
		class = nameArr[0]
	}
	return class
}

func (trMgr *roomManager) OneUserAnd3Robot() {

}
