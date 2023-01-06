package utils

import "strconv"

func GetStringFromInt64(i int64) string {
	return strconv.Itoa(int(i))
}
