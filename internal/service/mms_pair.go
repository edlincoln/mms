package service

import (
	"context"

	"github.com/edlincoln/mms/internal/dto"
	"github.com/edlincoln/mms/internal/http/request"
	"github.com/edlincoln/mms/internal/model"
	"github.com/edlincoln/mms/internal/repository"
	"github.com/edlincoln/mms/internal/utils"
	"github.com/edlincoln/mms/pkg/logger"
	"go.uber.org/zap"
)

const (
	period20  = 20
	period50  = 50
	period200 = 200
)

type MmsPairService interface {
	Save(ctx context.Context, pair model.MmsPairs) error
	BulkSave(ctx context.Context, pair string, candles []dto.Candle) error
	FindByPairAndTimestampRange(ctx context.Context, pair string, dto request.MmsPairDto) ([]model.MmsPairs, error)
	SmaCalculate(ctx context.Context, pair string, candles []dto.Candle) []model.MmsPairs
	Count(ctx context.Context, pair string) (int, error)
}

type MmsPairServiceImpl struct {
	repository repository.MmsPairRepository
}

func NewMmsPairService(repository repository.MmsPairRepository) MmsPairService {
	return MmsPairServiceImpl{repository: repository}
}

func (s MmsPairServiceImpl) FindByPairAndTimestampRange(ctx context.Context, pair string, dto request.MmsPairDto) ([]model.MmsPairs, error) {
	logger.Debug("finding mms pair")

	from, err := dto.GetFromAsDate()
	if err != nil {
		logger.Error("error getting from date", err)
		return nil, err
	}

	to, err := dto.GetToAsDate()
	if err != nil {
		logger.Error("error getting to date", err)
		return nil, err
	}

	return s.repository.FindByPairAndTimestampRange(ctx, pair, from, to)
}

func (s MmsPairServiceImpl) Save(ctx context.Context, pair model.MmsPairs) error {
	logger.Debug("inserting mms pair")
	err := s.repository.Save(ctx, pair)
	if err != nil {
		logger.Error("error getting to date", err)
		return err
	}
	return nil
}

func (s MmsPairServiceImpl) BulkSave(ctx context.Context, pair string, candles []dto.Candle) error {
	logger.Debug("bulk inserting mms pair")
	count, err := s.Count(ctx, pair)
	if err != nil {
		return err
	}

	if count > 0 {
		logger.Info("table already loaded.", zap.String("pair", pair))
		return nil
	}

	ret := s.SmaCalculate(ctx, pair, candles)
	return s.repository.BulkSave(ctx, ret)
}

func (s MmsPairServiceImpl) Count(ctx context.Context, pair string) (int, error) {
	logger.Debug("count pair records")
	count, err := s.repository.Count(ctx, pair)
	if err != nil {
		return count, err
	}
	return count, err
}

func (s MmsPairServiceImpl) SmaCalculate(ctx context.Context, pair string, candles []dto.Candle) []model.MmsPairs {
	result := make([]model.MmsPairs, 0)
	sum20 := float64(0)
	sum50 := float64(0)
	sum200 := float64(0)
	for i, value := range candles {
		count := i + 1
		sum20 += value.Close
		sum50 += value.Close
		sum200 += value.Close

		mms := new(model.MmsPairs)
		mms.Pair = pair
		mms.Timestamp = utils.UnixInt64ToDate(value.Timestamp)
		if i >= period20 {
			sum20 -= candles[i-period20].Close
			count = period20
			mms.Mms20 = sum20 / float64(count)
		} else {
			mms.Mms20 = 0
		}

		if i >= period50 {
			sum50 -= candles[i-period50].Close
			count = period50
			mms.Mms50 = sum50 / float64(count)
		} else {
			mms.Mms50 = 0
		}

		if i >= period200 {
			sum200 -= candles[i-period200].Close
			count = period200
			mms.Mms200 = sum200 / float64(count)
		} else {
			mms.Mms200 = 0
		}

		result = append(result, *mms)
	}

	return result
}
