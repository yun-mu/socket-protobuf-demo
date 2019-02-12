package newtime

import (
	"time"
)

// GetDayStartTimestampByClock 获取以 某时钟为基准的东八区的一天起始点时间戳
// 默认以24点为基准点
// eg: // 以 晚上 11 点为零点, clock = 23
func GetDayStartTimestampByClock(clock int) int64 {
	dayStartTimestamp := BeginningOfDay().UnixNano() / 1000000
	now := GetMilliTimestamp()

	if now >= (dayStartTimestamp + 3600*1000*(int64(clock-8))) {
		dayStartTimestamp += 3600 * 1000 * (int64(clock - 8))
	} else {
		dayStartTimestamp -= 3600 * 1000 * (int64(32 - clock))
	}
	return dayStartTimestamp
}

// GetWeekStartTimestampByClock 获取以 某时钟为基准的东八区的一周起始点时间戳
// 默认以24点为基准点
// eg: // 以 晚上 11 点为零点, clock = 23
func GetWeekStartTimestampByClock(clock int) int64 {
	weekStartTimestamp := BeginningOfWeek().UnixNano() / 1000000
	now := GetMilliTimestamp()

	if now >= (weekStartTimestamp + 3600*1000*160) {
		weekStartTimestamp += 3600 * 1000 * 160
	} else {
		weekStartTimestamp -= 3600 * 1000 * 8
	}
	return weekStartTimestamp - int64(24-clock)*3600*1000
}

// GetMilliTimestamp get the millisecond number since January 1, 1970 UTC
func GetMilliTimestamp() int64 {
	return time.Now().UnixNano() / 1000000
}
