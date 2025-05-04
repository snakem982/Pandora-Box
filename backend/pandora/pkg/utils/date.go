package utils

import (
	"time"
)

var layout = "2006-01-02 15:04"

// GetDateTime 获取当前时间并格式化为 "2006-01-02 15:04"
func GetDateTime() string {
	return time.Now().Format(layout)
}

// ParseDateTime 使用 time.Parse 将字符串解析为 time.Time 类型
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	parsedTime, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
