// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package kitchen

import (
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/app"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/handler/consumers"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/cache"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq/consumer"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq/producer"
)

// Injectors from wire.go:

func Initialize(filePath configs.ConfigFilePath) (*app.App, func(), error) {
	config, err := configs.GetConfigFromYaml(filePath)
	if err != nil {
		return nil, nil, err
	}
	client := producer.NewClient(config)
	consumerConsumer := consumer.NewConsumer(config)
	consumerHandler := consumers.NewKafkaHandler(consumerConsumer)
	redis := config.Redis
	cacheClient := cache.NewRedisClient(redis)
	kitchenService := service.NewService(cacheClient)
	appApp := app.NewApp(client, consumerHandler, kitchenService)
	return appApp, func() {
	}, nil
}
