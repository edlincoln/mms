package configs

import (
	"context"
	"time"

	"github.com/edlincoln/mms/internal/container"
	"github.com/edlincoln/mms/internal/global"
	"github.com/edlincoln/mms/pkg/logger"
)

type DataExtractor struct{}

func (d *DataExtractor) Init(ctx context.Context) error {

	if !global.AppConfig.Extractor.Enabled {
		logger.Info("extractor disabled. Skiping data load.")
		return nil
	}

	logger.Info("extractor enabled. Starting data load.")

	url := global.AppConfig.Extractor.Url

	r := global.AppConfig.Extractor.Range
	pairs := global.AppConfig.Extractor.Pairs

	to := time.Now().AddDate(0, 0, -1)
	from := to.AddDate(0, 0, -r)

	err := container.Container().ExtractorManagerService.LoadData(ctx, url, pairs, from.Unix(), to.Unix())
	if err != nil {
		logger.Error("error extracting data, try again latter.", err)
		return err
	}

	logger.Info("data load fineshed.")

	return nil
}

func (d *DataExtractor) Close(ctx context.Context) error {
	return nil
}
