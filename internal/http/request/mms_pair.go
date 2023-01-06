package request

import (
	"time"

	"github.com/edlincoln/mms/internal/http/handlers"
	"github.com/edlincoln/mms/internal/utils"
)

const daysAllowed int = 365

type MmsPairDto struct {
	From  string  `json:"from"  validate:"required" example:"1577836800"`
	To    *string `json:"to"    example:"1577836800"`
	Range string  `json:"range" validate:"required,oneof=20 50 100" example:"20"`
}

func (dto MmsPairDto) Valid() handlers.RestErr {
	limit := time.Now().AddDate(0, 0, -daysAllowed)
	dtFrom, _ := utils.UnixStrToDate(dto.From)
	if dtFrom.Before(limit) {
		return handlers.NewBadRequestError("data inicio anterior a data limite")
	}
	return nil
}

func (dto MmsPairDto) GetFromAsDate() (time.Time, error) {
	return utils.UnixStrToDate(dto.From)
}

func (dto MmsPairDto) GetToAsDate() (time.Time, error) {
	if dto.To == nil {
		return time.Now().AddDate(0, 0, -1), nil
	}

	return utils.UnixStrToDate(*dto.To)
}
