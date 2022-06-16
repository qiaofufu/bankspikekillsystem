package utils

import (
	"errors"
	"time"
)

var LOC, _ = time.LoadLocation("Asia/Shanghai")

// GetTime 获取time
func GetTime(timeStr string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"

	ret, err := time.ParseInLocation(layout, timeStr, LOC)

	if err != nil {
		return time.Now(), err
	} else {
		return ret, nil
	}
}

// GetDuration 获取时间间隔
func GetDuration(durationStr string) (duration time.Duration, err error) {
	duration, err = time.ParseDuration(durationStr)
	if err != nil {
		err = errors.New("获取时间间隔失败")
	}
	return
}
