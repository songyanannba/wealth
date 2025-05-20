package websocket

// 5斑马 4犀牛  3大象 2狮子  1蛇  1皇冠   1 LUCKY  1 大猩猩

type AnimalConfig struct {
	Seat     int
	AnimalId int
}

// GetAnimalWheel 获取最外圈的动物固定排序
// 1 大猩猩    1
// 2 LUCKY    1
// 3 皇冠      1  //皇冠 代表  狮子 大象  犀牛  斑马
// 4 蛇       1
// 5 狮子     2
// 6 大象     3
// 7 犀牛     4
// 8 斑马     5
func GetAnimalWheel() []*AnimalConfig {
	var (
		animalConfigs = make([]*AnimalConfig, 18)
	)

	//animalId := []int{1, 2, 3, 4, 5, 5, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 8, 8}
	//根据微信发的手动排序
	animalId := []int{8, 6, 7, 1, 8, 2, 7, 8, 6, 5, 7, 8, 4, 6, 8, 7, 3, 5}

	for i, v := range animalId {
		animalConfigs[i] = &AnimalConfig{Seat: i, AnimalId: v}
	}
	//helper.SliceShuffle(animalConfigs)
	return animalConfigs
}
