package common

// AppId 1：未定义 2: 钓鱼 3：酒馆故事
const AppId1 = 1         //默认
const AppId = 2          //钓鱼
const AppId3 = 3         //酒馆故事
const AppId4 = 4         //幽影魔塔
const AppId10 = 10       //meme
const AwaitNextTime = 20 //

const MagicTowerId = "meme_battle"

//const ApiLimit = true

// ConvertFLevel 获取鱼等级
func ConvertFLevel(fLevel int) string {
	var res string
	res = "Z"
	if fLevel == 1 {
		res = "A"
	} else if fLevel == 2 {
		res = "B"
	} else if fLevel == 3 {
		res = "C"
	} else if fLevel == 4 {
		res = "D"
	} else if fLevel == 5 {
		res = "E"
	} else if fLevel == 6 {
		res = "L"
	}
	return res
}

func FLevelAdopt(fLevel int) float64 {
	var res float64
	//todo
	//第一版的预期是，全部300，a200，b100，c50，d20
	if fLevel == 1 {
		res = 200
	} else if fLevel == 2 {
		res = 100
	} else if fLevel == 3 {
		res = 50
	} else if fLevel == 4 {
		res = 20
	} else if fLevel == 0 {
		res = 300
	}
	return res
}
