package utils

import (
	"time"
)

// 获取系统本地时区
var localZone, _ = time.LoadLocation("Local")

var layout = "2006-01-02 15:04"

// GetDateTime 获取当前时间，并格式化为 "2006-01-02 15:04"
func GetDateTime() string {
	return time.Now().In(localZone).Format(layout)
}

// ParseDateTime 使用 time.ParseInLocation 将字符串解析为 time.Time 类型，并设定本地时区
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	parsedTime, err := time.ParseInLocation(layout, dateTimeStr, localZone)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

// GetHourDiff 计算当前时间与输入时间的时间差（小时数），返回 int 类型
func GetHourDiff(inputTime time.Time) int {
	currentTime := time.Now().In(localZone) // 获取本地时区的当前时间
	duration := currentTime.Sub(inputTime)  // 计算时间差

	return int(duration.Hours()) // 返回 int 类型的小时数
}
