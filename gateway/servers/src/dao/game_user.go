package dao

import (
	"gateway/global"
	"gateway/models/table"
)

func SaveGameUser(val *table.GameUser) error {
	global.GVA_LOG.Infof("SaveGameUser%v", *val)
	user, err := table.GetGameUserByUid(val.UserId)
	//if err != nil {
	//	return err
	//}
	if user.ID > 0 {
		//修改
		if user.Nickname != val.Nickname {
			user.Nickname = val.Nickname
			user.KingCoin = val.KingCoin
			err = table.SaveGameUser(user)
		}
	} else {
		err = table.CreateGameUser(val)
	}
	return err
}

func GetGameUser(uid string) (val *table.GameUser, err error) {
	return table.GetGameUserByUid(uid)
}
