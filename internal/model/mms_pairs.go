package model

import "time"

type MmsPairs struct {
	Id        int64
	Pair      string
	Timestamp time.Time
	Mms20     float64
	Mms50     float64
	Mms200    float64
}
