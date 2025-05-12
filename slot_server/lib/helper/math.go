package helper

import (
	"github.com/shopspring/decimal"
)

func Sum[T Number](num ...T) float64 {
	if len(num) == 0 {
		return 0
	}
	var sum decimal.Decimal
	for _, v := range num {
		sum = sum.Add(decimal.NewFromFloat(float64(v)))
	}
	f, _ := sum.Float64()
	return f
}

func Mul[T Number, V Number](d1 T, num ...V) float64 {
	if len(num) == 0 || d1 == 0 {
		return 0
	}
	var sum = decimal.NewFromFloat(float64(d1))
	for _, v := range num {
		if v == 0 {
			return 0
		}
		sum = sum.Mul(decimal.NewFromFloat(float64(v)))
	}
	f, _ := sum.Float64()
	return f
}

func Div[T Number, V Number](d1 T, num ...V) float64 {
	if len(num) == 0 || d1 == 0 {
		return 0
	}
	var sum = decimal.NewFromFloat(float64(d1))
	for _, v := range num {
		if v == 0 {
			return 0
		}
		sum = sum.Div(decimal.NewFromFloat(float64(v)))
	}
	f, _ := sum.Float64()
	return f
}

func MulToInt[T Number, V Number](d1 T, num ...V) int64 {
	if len(num) == 0 || d1 == 0 {
		return 0
	}
	var sum = decimal.NewFromFloat(float64(d1))
	for _, v := range num {
		if v == 0 {
			return 0
		}
		sum = sum.Mul(decimal.NewFromFloat(float64(v)))
	}
	return sum.IntPart()
}

func Abs[T Signed](num T) T {
	if num < 0 {
		return -num
	}
	return num
}

// NearKey 匹配最近值的key nums为降序
func NearKey[T Signed](nums []T, num T) []int {
	minDiff := nums[0] - num
	keys := []int{0}
	for i := 1; i < len(nums); i++ {
		diff := Abs(nums[i] - num)
		if diff < minDiff {
			minDiff = diff
			keys = []int{i}
		} else if diff == minDiff {
			keys = append(keys, i)
		} else {
			break
		}
	}
	return keys
}

// NearVal 匹配最近值 nums为降序
func NearVal[T Signed](nums []T, num T) T {
	minDiff := nums[0] - num
	for i := 1; i < len(nums); i++ {
		diff := Abs(nums[i] - num)
		if diff > minDiff {
			return nums[i-1]
		} else {
			minDiff = diff
		}
	}
	return nums[len(nums)-1]
}

func Mul100(num float64) int64 {
	return decimal.NewFromFloat(num).Mul(decimal.NewFromInt(100)).IntPart()
}

func Div100(num int64) float64 {
	v, _ := decimal.NewFromInt(num).Div(decimal.NewFromInt(100)).Float64()
	return v
}

func FlotDiv100(num float64) float64 {
	v, _ := decimal.NewFromFloat(num).Div(decimal.NewFromInt(100)).Float64()
	return v
}

func Range(n1, n2 int) []int {
	numbers := make([]int, n2-n1+1)

	for i := range numbers {
		numbers[i] = n1 + i
	}
	return numbers
}

func CeilDiv(a, b int) int {
	q, r := a/b, a%b
	if r == 0 {
		return q
	}
	if (r > 0 && b > 0) || (r < 0 && b < 0) {
		q++
	}
	return q
}
