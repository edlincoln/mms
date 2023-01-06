package service

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/edlincoln/mms/internal/dto"
	"github.com/edlincoln/mms/internal/http/request"
	"github.com/edlincoln/mms/internal/model"
	"github.com/edlincoln/mms/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	zero   float64 = 0
	mms20  float64 = 5
	mms50  float64 = 5.5
	mms200 float64 = 6
)

func TestFindByPairAndTimestampRange(t *testing.T) {
	initValues()
	t.Run("load data", func(t *testing.T) {
		ctx := context.Background()
		repoMock := new(repository.MockMmsPairRepository)
		repoMock.On("FindByPairAndTimestampRange", ctx, brlbtc, fromDt, toDt).Return(GetMmsPairs(), nil)

		mmsService := NewMmsPairService(repoMock)
		ret, err := mmsService.FindByPairAndTimestampRange(ctx, brlbtc, GetMmsPairDto())
		assert.Nil(t, err)
		assert.NotNil(t, ret)
		assert.Equal(t, 2, len(ret))
	})

	t.Run("load data error", func(t *testing.T) {
		ctx := context.Background()
		repoMock := new(repository.MockMmsPairRepository)
		repoMock.On("FindByPairAndTimestampRange", ctx, brlbtc, fromDt, toDt).Return(nil, errors.New("error finding"))

		mmsService := NewMmsPairService(repoMock)
		ret, err := mmsService.FindByPairAndTimestampRange(ctx, brlbtc, GetMmsPairDto())
		assert.NotNil(t, err)
		assert.Nil(t, ret)
	})
}

func TestSaveMmsPair(t *testing.T) {
	initValues()
	t.Run("save pair", func(t *testing.T) {
		ctx := context.Background()
		repoMock := new(repository.MockMmsPairRepository)
		repoMock.On("Save", ctx, GetMmsPair()).Return(nil)

		mmsService := NewMmsPairService(repoMock)
		err := mmsService.Save(ctx, GetMmsPair())
		assert.Nil(t, err)
	})

	t.Run("save pair", func(t *testing.T) {
		ctx := context.Background()
		repoMock := new(repository.MockMmsPairRepository)
		repoMock.On("Save", ctx, GetMmsPair()).Return(errors.New("error"))

		mmsService := NewMmsPairService(repoMock)
		err := mmsService.Save(ctx, GetMmsPair())
		assert.NotNil(t, err)
	})
}

func TestBulkSaveMmsPair(t *testing.T) {
	initValues()
	t.Run("save pair", func(t *testing.T) {
		ctx := context.Background()
		repoMock := new(repository.MockMmsPairRepository)
		repoMock.On("Count", ctx, brlbtc).Return(0, nil)
		repoMock.On("BulkSave", ctx, mock.Anything).Return(nil)

		mmsService := NewMmsPairService(repoMock)

		err := mmsService.BulkSave(ctx, brlbtc, GetCandles())
		assert.Nil(t, err)
	})

	t.Run("save pair", func(t *testing.T) {
		ctx := context.Background()
		repoMock := new(repository.MockMmsPairRepository)
		repoMock.On("Count", ctx, brlbtc).Return(0, errors.New("error counting"))

		mmsService := NewMmsPairService(repoMock)

		err := mmsService.BulkSave(ctx, brlbtc, GetCandles())
		assert.NotNil(t, err)
	})

	t.Run("save pair", func(t *testing.T) {
		ctx := context.Background()
		repoMock := new(repository.MockMmsPairRepository)
		repoMock.On("Count", ctx, brlbtc).Return(1, nil)

		mmsService := NewMmsPairService(repoMock)

		err := mmsService.BulkSave(ctx, brlbtc, GetCandles())
		assert.Nil(t, err)
	})
}

func GetMmsPairDto() request.MmsPairDto {
	to := strconv.Itoa(int(to))
	return request.MmsPairDto{
		From:  strconv.Itoa(int(from)),
		To:    &to,
		Range: "20",
	}
}

func GetMmsPairs() []model.MmsPairs {
	list := make([]model.MmsPairs, 0)
	list = append(list, model.MmsPairs{Mms20: mms20, Mms50: zero, Mms200: zero})
	list = append(list, model.MmsPairs{Mms20: mms20, Mms50: mms50, Mms200: mms200})

	return list
}

func GetMmsPair() model.MmsPairs {
	return model.MmsPairs{Mms20: mms20, Mms50: mms50, Mms200: mms200}
}

func GetCandles() []dto.Candle {
	list := make([]dto.Candle, 0)
	list = append(list, dto.Candle{Timestamp: from, Close: mms20})
	list = append(list, dto.Candle{Timestamp: from, Close: mms50})
	list = append(list, dto.Candle{Timestamp: from, Close: mms200})
	return list
}
