package helper

import (
	"slot_server/lib/global"
	"time"
)

func YearMonthDayStr() string {
	// 获取当前时间z
	now := time.Now()
	return now.Format("2006-01-02")
}

//func YearMonthStr() string {
//	// 获取当前时间
//	now := time.Now()
//	return now.Format("2006-01")
//}

func TimeIntToStr(t int64) string {
	// 获取当前时间
	unix := time.Unix(t, 0)
	return unix.Format("2006-01-02 15:04:05")
}

func PreYearMonthDayStr() string {
	// 获取当前时间
	currentTime := time.Now()
	preTime := currentTime.AddDate(0, 0, -1) //前5天就写-5。
	return preTime.Format("2006-01-02")
}

func LocalTime() *time.Time {
	// 设置时区为东八区
	//utcTime := time.Now().UTC()
	localTime := time.Now().In(global.Location)
	return &localTime
}
