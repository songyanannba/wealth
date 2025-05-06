package helper

import "strconv"

func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Itoa[V Int](i V) string {
	return strconv.Itoa(int(i))
}

type Int interface {
	int8 | int | int32 | int64 | uint8 | uint16 | uint32 | uint | uint64
}

type Number interface {
	Int | float32 | float64
}

type Signed interface {
	int8 | int32 | int | int64 | float64
}

type Unsigned interface {
	uint8 | uint16 | uint32 | uint | uint64
}
