package helper

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"time"
)

// RandArr 数组中随机一个值
func RandArr[T any](arr []T) T {
	var res T
	if len(arr) == 0 {
		return res
	}
	bigInt, _ := crand.Int(crand.Reader, big.NewInt(int64(len(arr))))
	return arr[bigInt.Int64()]
}

// RandInt 随机一个值
func RandInt(v int) int {
	if v <= 0 {
		return 0
	}
	bigInt, _ := crand.Int(crand.Reader, big.NewInt(int64(v)))
	return int(bigInt.Int64())
}
func GetRand(n int) int {
	rand.Seed(time.Now().UnixNano())
	// 获取一个随机整数
	randomNumber := rand.Intn(n) // 生成一个0到n的随机整数
	return randomNumber
}

func RandGetValue[T any](ts []T) T {
	if len(ts) == 0 {
		return *new(T)
	}
	return ts[RandInt(len(ts))]
}

// RandScope 随意一个范围内的值
func RandScope(s, e int) int {
	if e <= s {
		return s
	}
	return RandInt(e-s+1) + s
}

// RandomLongWeight 根据权重随机一个值
func RandomLongWeight(weights []int) int {
	length := len(weights)
	if length <= 1 {
		if length == 1 {
			return 1
		}
		return -1
	}
	start := weights[0]
	end := weights[length-1]

	res := RandScope(start, end-1)
	return InArrIntervalIndex(weights, res)
}

func InArrIntervalIndex(arr []int, num int) int {
	for i := len(arr) - 2; i >= 0; i-- {
		if num >= arr[i] {
			return i
		}
	}
	return 0
}
