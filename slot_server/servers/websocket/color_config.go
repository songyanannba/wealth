package websocket

import "slot_server/lib/helper"

//3黄+7绿+7红   +三色
//4黄+7绿+6红   +三色
//4黄+6绿+7红
//5黄+5绿+7红
//5黄+7绿+5红
//5黄+6绿+6红
//6黄+6绿+5红
//6黄+5绿+6红
//6黄+7绿+4红
//6黄+4绿+7红
//7黄+5绿+5红
//7黄+6绿+4红
//7黄+4绿+6红
//7黄+7绿+3红
//7黄+3绿+7红   +三色

type ColorConfig struct {
	Seat    int
	ColorId int
}

func GetColorAllClass() [][]int {
	colorAllClass := [][]int{}

	// 1 黄
	// 2 绿
	// 3 红
	// 4 三色

	//3黄+7绿+7红   +三色
	color1 := []int{1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 4}
	colorAllClass = append(colorAllClass, color1)

	//4黄+7绿+6红   +三色
	color2 := []int{1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 4}
	colorAllClass = append(colorAllClass, color2)

	////4黄+6绿+7红
	color3 := []int{1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 4}
	colorAllClass = append(colorAllClass, color3)

	return colorAllClass
}

func GetColorWheel() []*ColorConfig {
	colorConfigs := make([]*ColorConfig, 18)
	getColorAllClass := GetColorAllClass()
	helper.RandInt(len(getColorAllClass))
	colorAllClass := getColorAllClass[helper.RandInt(len(getColorAllClass))]
	for i, v := range colorAllClass {
		colorConfigs[i] = &ColorConfig{Seat: i, ColorId: v}
	}
	helper.SliceShuffle(colorConfigs)
	return colorConfigs
}
