package response

import (
	"github.com/edlincoln/mms/internal/model"
	"github.com/edlincoln/mms/internal/utils"
)

type MmsPair struct {
	Timestamp int64   `json:"timestamp"`
	Mms       float64 `json:"mms"`
}

func ToMmsPairListResponse(r string, pairs []model.MmsPairs) []MmsPair {
	ret := make([]MmsPair, 0)
	for _, p := range pairs {
		switch r {
		case "20":
			ret = append(ret, getMmsPair(p.Timestamp.Unix(), p.Mms20))
		case "50":
			ret = append(ret, getMmsPair(p.Timestamp.Unix(), p.Mms50))
		case "200":
			ret = append(ret, getMmsPair(p.Timestamp.Unix(), p.Mms200))
		}
	}
	return ret
}

func getMmsPair(t int64, m float64) MmsPair {
	return MmsPair{Timestamp: t, Mms: utils.GetRoundedFloatValue(m)}
}
