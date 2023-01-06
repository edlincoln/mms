package container

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/edlincoln/mms/internal/global"
	"github.com/edlincoln/mms/internal/http/client"
	"github.com/edlincoln/mms/internal/repository"
	"github.com/edlincoln/mms/internal/service"
	"github.com/edlincoln/mms/pkg/logger"
	"github.com/go-co-op/gocron"
)

var container *ServiceContainer
var once sync.Once

type ServiceContainer struct {
	DatabaseManager         repository.DatabaseManager
	MmsPairService          service.MmsPairService
	ExtractorManagerService service.ExtractorManagerService
}

func Container() *ServiceContainer {
	once.Do(func() {
		databaseManager := repository.NewDatabaseManager()
		mmsPairService := initMmsPairService(databaseManager)
		extractorManagerService := initExtractorManagerService(mmsPairService)
		container = &ServiceContainer{
			DatabaseManager:         databaseManager,
			MmsPairService:          mmsPairService,
			ExtractorManagerService: extractorManagerService,
		}
		initDailyScheduler()
	})
	return container
}

func initDailyScheduler() {
	logger.Info("starting scheduler")
	sched := gocron.NewScheduler(time.Local)

	retry := global.AppConfig.Extractor.Daily.Retry

	// sched.Cron("*/1 * * * *").Do(func() {
	// 	Container().ExtractorManagerService.LoadDailyData(context.Background(), retry)
	// })

	days := global.AppConfig.Extractor.Daily.Frequency
	hour := global.AppConfig.Extractor.Daily.Hour
	sched.Every(days).Day().At(hour).Do(func() {
		Container().ExtractorManagerService.LoadDailyData(context.Background(), retry)
	})

	sched.StartAsync()
}

func initMmsPairService(databaseManager repository.DatabaseManager) service.MmsPairService {
	repo := repository.NewMmsPairRepository(&databaseManager)
	return service.NewMmsPairService(repo)
}

func initExtractorManagerService(mmsPairService service.MmsPairService) service.ExtractorManagerService {
	httpClient := client.NewHttpClient(&http.Client{})
	return service.NewExtractorManagerService(mmsPairService, httpClient)
}
