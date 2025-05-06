package helper

import (
	"fmt"
	"testing"
)

func Test_rand(t *testing.T) {
	for i := 0; i < 100; i++ {
		randInt := RandInt(3)
		//0 1 2
		fmt.Println(randInt)
	}
}
