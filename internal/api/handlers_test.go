package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"news_aggregator/internal/models"
	"news_aggregator/internal/service"
	"news_aggregator/internal/storage/postgres"
	"os"
	"testing"
)

func TestHandler_GetNews(t *testing.T) {
	config := &models.DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
		Name:     getEnvOrDefault("DB_NAME", "news_aggregator_test"),
		SSLMode:  "disable",
	}

	storage, err := postgres.New(config)
	if err != nil {
		t.Skip("Пропуск теста: нет подключения к тестовой БД")
	}
	defer func(storage *postgres.Storage) {
		err := storage.Close()
		if err != nil {

		}
	}(storage)

	newsService := service.NewNewsService(storage, []string{"https://test.com"})
	handler := NewHandler(newsService)

	req := httptest.NewRequest("GET", "/api/news/10", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/news/{n}", handler.GetNews)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusOK, w.Code)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
