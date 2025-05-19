package websocket

import "slot_server/lib/helper"

type SizeConfig struct {
	Seat int
}

//一共16个

//1 星星+三色 180

type BetZoneConfig struct {
	Seat     int
	AnimalId int
	ColorId  []int
	Size     int //1大(粉色) 2小（紫色）
	BetRate  float64
}

func GetBetZoneFigure() []*BetZoneConfig {
	betZoneConfigs := make([]*BetZoneConfig, 16)

	// 1 大猩猩    1
	// 2 LUCKY    1
	// 3 皇冠      1
	// 4 蛇        1
	// 5 狮子      2
	// 6 大象     3
	// 7 犀牛     4
	// 8 斑马     5

	// 1 黄
	// 2 绿
	// 3 红

	for i := 0; i < 16; i++ {
		var (
			animalId = 0
			size     = 0
			betRate  float64
			colorId  = make([]int, 0)
		)

		if i == 0 {
			animalId = 1
			colorId = []int{1, 2, 3, 4}
			betRate = 180
		}

		if i == 1 {
			animalId = 5
			colorId = []int{3}
			betRate = 46

		}
		if i == 2 {
			animalId = 5
			colorId = []int{2}
			betRate = 40
		}
		if i == 3 {
			animalId = 5
			colorId = []int{1}
			betRate = 25
		}
		//蛇
		if i == 4 {
			animalId = 4
			colorId = []int{1, 2, 3, 4}
			betRate = 120

		}
		//
		if i == 5 {
			animalId = 6
			colorId = []int{3}
			betRate = 23

		}
		if i == 6 {
			animalId = 6
			betRate = 20
			colorId = []int{2}

		}
		if i == 7 {
			animalId = 6
			betRate = 12
			colorId = []int{1}
		}
		//
		if i == 8 {
			size = 1
			betRate = 2
		}
		if i == 9 {
			animalId = 7
			colorId = []int{3}
			betRate = 13
		}
		if i == 10 {
			animalId = 7
			colorId = []int{2}
			betRate = 11

		}
		if i == 11 {
			animalId = 7
			colorId = []int{1}
			betRate = 7
		}
		if i == 12 {
			size = 2
			betRate = 2
		}
		if i == 13 {
			animalId = 8
			colorId = []int{3}
			betRate = 8
		}
		if i == 14 {
			animalId = 8
			colorId = []int{2}
			betRate = 7
		}
		if i == 15 {
			animalId = 8
			colorId = []int{1}
			betRate = 4
		}

		betZoneConfigs[i] = &BetZoneConfig{
			Seat:     i,
			AnimalId: animalId,
			ColorId:  colorId,
			Size:     size,
			BetRate:  betRate,
		}
	}

	return betZoneConfigs
}

func GetBetZoneConfigByAnimalIdAndColorId(animalId, colorId int) *BetZoneConfig {
	res := &BetZoneConfig{}

	for _, betZoneFigure := range GetBetZoneFigure() {
		if !helper.InArr(colorId, betZoneFigure.ColorId) {
			continue
		}
		if betZoneFigure.AnimalId == animalId {
			res = betZoneFigure
			break
		}
	}
	return res
}
