// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package batch

import (
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/app"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/handler/jobs"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/repo"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/cache"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/data_access"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/s3"
)

// Injectors from wire.go:

func Initialize(filePath configs.ConfigFilePath) (*app.App, func(), error) {
	config, err := configs.GetConfigFromYaml(filePath)
	if err != nil {
		return nil, nil, err
	}
	configsS3 := config.S3
	s3Client := s3.NewS3Client(configsS3)
	db, cleanup, err := database.InitializeAndMigrateUpDB(config)
	if err != nil {
		return nil, nil, err
	}
	goquDatabase := database.InitializeGoquDB(db)
	orderRepo := repo.NewOrderRepo(goquDatabase)
	redis := config.Redis
	client := cache.NewRedisClient(redis)
	batchService := service.NewBatchService(orderRepo, client)
	clearAllOrderEOD := jobs.NewExecuteClearAllOrderEOD(batchService)
	cron := config.Cron
	appApp := app.NewApp(s3Client, batchService, clearAllOrderEOD, cron)
	return appApp, func() {
		cleanup()
	}, nil
}