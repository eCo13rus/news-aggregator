package main

import (
	"encoding/json"
	"log"
	"news_aggregator/internal/api"
	"news_aggregator/internal/models"
	"news_aggregator/internal/service"
	"news_aggregator/internal/storage/postgres"
	"os"
	"time"
)

func main() {
	config, err := loadConfig("configs/config.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	storage, err := postgres.New(&config.Database)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer func(storage *postgres.Storage) {
		err := storage.Close()
		if err != nil {

		}
	}(storage)

	newsService := service.NewNewsService(storage, config.RSS)

	period := time.Duration(config.Period) * time.Minute
	if err := newsService.Start(period); err != nil {
		log.Fatalf("Ошибка запуска сервиса новостей: %v", err)
	}
	defer newsService.Stop()

	handler := api.NewHandler(newsService)

	server := api.NewServer(handler)
	server.SetupRoutes()

	addr := config.Server.Host + ":" + config.Server.Port
	if err := server.Start(addr); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func loadConfig(path string) (*models.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var config models.Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
