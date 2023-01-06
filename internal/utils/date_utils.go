package utils

import (
	"strconv"
	"time"
)

const (
	apiDateLayout = "2006-01-02T15:04:05Z"
	apiDbLayout   = "2006-01-02 15:04:05"
)

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}

func GetNowDBFormat() string {
	return GetNow().Format(apiDbLayout)
}

func UnixStrToDate(ts string) (time.Time, error) {
	dt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(dt, 0), nil
}

func UnixInt64ToDate(ts int64) time.Time {
	return time.Unix(ts, 0)
}
