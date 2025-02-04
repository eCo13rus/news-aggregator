package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"news_aggregator/internal/models"
	"news_aggregator/internal/service"
	"news_aggregator/internal/storage/postgres"
	"testing"
)

func TestHandler_GetNews(t *testing.T) {
	config := &models.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "eco",
		Password: "Horror4199.",
		Name:     "news_aggregator",
		SSLMode:  "disable",
	}

	storage, err := postgres.New(config)
	if err != nil {
		t.Skip("Пропуск теста: нет подключения к БД")
	}
	defer storage.Close()

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
