package utils

import "math"

func GetRoundedFloatValue(value float64) float64 {
	return math.Round(value*100) / 100
}
