package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/edlincoln/mms/internal/dto"
	"github.com/edlincoln/mms/internal/global"
	"github.com/edlincoln/mms/internal/http/client"
	"github.com/edlincoln/mms/pkg/logger"
)

type ExtractorManagerService interface {
	LoadData(ctx context.Context, url string, pairs []string, from, to int64) error
	LoadDailyData(ctx context.Context, retry int) error
}

type ExtractorManagerServiceImp struct {
	mmsPairService MmsPairService
	httpClient     client.HttpClient
}

func NewExtractorManagerService(mmsPairService MmsPairService, httpClient client.HttpClient) ExtractorManagerService {
	return ExtractorManagerServiceImp{mmsPairService: mmsPairService, httpClient: httpClient}
}

func (s ExtractorManagerServiceImp) LoadData(ctx context.Context, url string, pairs []string, from, to int64) error {

	logger.Debug("extracting data")

	if global.AppConfig.Extractor.Mocked {
		eth := s.mockeExtractorResponse("/resources/mock/eth.json")
		s.mmsPairService.BulkSave(ctx, "BRLETH", eth)

		btc := s.mockeExtractorResponse("/resources/mock/btc.json")
		s.mmsPairService.BulkSave(ctx, "BRLBTC", btc)

		return nil
	}

	for _, v := range pairs {
		count, err := s.mmsPairService.Count(ctx, v)
		if err != nil {
			logger.Error("error extracting data", err)
			return err
		}
		if count == 0 {
			candles, err := s.extractData(ctx, url, v, from, to)
			if err != nil {
				logger.Error("error extracting data", err)
				return err
			}

			err = s.mmsPairService.BulkSave(ctx, v, candles)
			if err != nil {
				logger.Error("error extracting data", err)
				return err
			}
		}
	}
	return nil
}

func (s ExtractorManagerServiceImp) LoadDailyData(ctx context.Context, retry int) error {
	logger.Info("extracting daily data")

	if global.AppConfig.Extractor.Mocked {
		eth := s.mockeExtractorResponse("/resources/mock/eth.json")
		s.mmsPairService.BulkSave(ctx, "BRLETH", eth)

		btc := s.mockeExtractorResponse("/resources/mock/btc.json")
		s.mmsPairService.BulkSave(ctx, "BRLBTC", btc)

		return nil
	}

	period := global.AppConfig.Extractor.Daily.MaxPeriod
	retryS := global.AppConfig.Extractor.Daily.RetryTime
	pairs := global.AppConfig.Extractor.Pairs
	url := global.AppConfig.Extractor.Url

	to := time.Now().AddDate(0, 0, -1)
	from := to.AddDate(0, 0, -period)

	for _, v := range pairs {
		candles, err := s.extractData(ctx, url, v, from.Unix(), to.Unix())
		if err != nil {
			if retry > 0 {
				retry--
				logger.Debug(fmt.Sprintf("error extracting daily data. retry: %d", retry))
				time.Sleep(retryS)
				s.LoadDailyData(ctx, retry)
			}
			logger.Error("error extracting data", err)
			return err
		}

		smas := s.mmsPairService.SmaCalculate(ctx, v, candles)
		lastSma := smas[len(smas)-1:][0]
		err = s.mmsPairService.Save(ctx, lastSma)
		if err != nil {
			if retry > 0 {
				retry--
				logger.Debug(fmt.Sprintf("error extracting daily data. retry: %d", retry))
				time.Sleep(retryS)
				s.LoadDailyData(ctx, retry)
			}
			logger.Error("error extracting data", err)
			return err
		}
	}
	return nil
}

func (s *ExtractorManagerServiceImp) extractData(ctx context.Context, url, pair string, from, to int64) ([]dto.Candle, error) {
	var params = map[string]string{"from": strconv.Itoa(int(from)), "to": strconv.Itoa(int(to)), "precision": "1d"}
	res, err := s.httpClient.Get(ctx, fmt.Sprintf(url, pair), nil, params)
	if err != nil {
		logger.Error("error extracting data", err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		logger.Debug(res.Status)
		return nil, errors.New(res.Status)
	}

	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		logger.Error("error extracting data", err)
		return nil, err
	}

	var responseObject dto.ExtractorResponse
	err = json.Unmarshal(body, &responseObject)
	if err != nil {
		logger.Error("error extracting data", err)
		return nil, err
	}

	return responseObject.Candles, nil
}

func (s *ExtractorManagerServiceImp) mockeExtractorResponse(filePath string) []dto.Candle {
	dirname, err := os.Getwd()

	fmt.Println("!!!!" + dirname + "!!!!")

	if err != nil {
		panic(err)
	}

	// body, err := os.ReadFile(path.Join(dirname, "../") + filePath)
	body, err := os.ReadFile(path.Join(dirname,  filePath))
	if err != nil {
		log.Fatal(err)
	}

	var responseObject dto.ExtractorResponse
	json.Unmarshal(body, &responseObject)

	return responseObject.Candles
}
