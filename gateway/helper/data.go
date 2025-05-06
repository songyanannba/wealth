package helper

import "time"

func YearMonthDayStr() string {
	// 获取当前时间
	now := time.Now()
	return now.Format("2006-01-02")
}

func YearMonthStr() string {
	// 获取当前时间
	now := time.Now()
	return now.Format("2006-01")
}

func TimeIntToStr(t int64) string {
	// 获取当前时间
	unix := time.Unix(t, 0)
	return unix.Format("2006-01-02 15:04:05")
}

func LocalTime() *time.Time {
	utcTime := time.Now().UTC()
	localTime := utcTime.In(time.Local)
	return &localTime
}
