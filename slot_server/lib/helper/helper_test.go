package helper

import (
	"fmt"
	"tavern_story_room/global"
	"testing"
	"time"
)

func Test_Marshal(t *testing.T) {
	//age18 := AttenuationByAge(18)
	//fmt.Println("age18", age18)
	//age19 := AttenuationByAge(19)
	//fmt.Println("age19", age19)
	//attenuationRandInt, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", age19), 64)
	//fmt.Println("attenuationRandInt", attenuationRandInt)
	//fmt.Println("age19", int(age19))
	//age20 := AttenuationByAge(20)
	//fmt.Println("age20", age20)
	var num float64
	for i := 18; i < 38; i++ {
		//age := AttenuationByAge(i)
		age := NewAttenuationByAge(i, 5, 0.08)
		fmt.Println("i == ", i, age, Mul(Mul(age, 24), 0.08))
		num += Mul(Mul(age, 24), 0.08)
	}
	fmt.Println("num == ", num)

}

func Test_LocalTime(t *testing.T) {

	global.Location, _ = time.LoadLocation("Asia/Shanghai")

	time := LocalTime()

	fmt.Println("time == ", time)
}
