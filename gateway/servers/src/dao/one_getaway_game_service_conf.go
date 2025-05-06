package dao

import (
	"gateway/models/table"
)

func GetGameServiceConf(id int) (val *table.GameServiceConf, err error) {
	return table.GetGameServiceConf(id)
}
